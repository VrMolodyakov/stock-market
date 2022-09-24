import { useParams } from "react-router-dom";
import useAuth from '../routing/useAuth';
import { useState, useEffect, useMemo } from 'react'
import Chart from 'react-apexcharts';
import axios from "axios";
import jwt_decode from 'jwt-decode'
import "./StockCode.css"

const Code = (props) => {
  const { slug } = useParams();
  const [stockData, setStockData] = useState({});
  const [priceTime, setPriceTime] = useState(null);
  const [symbol, setSymbol] = useState("");
  const [priceInfo, setPriceInfo] = useState([{
    data: []
  }]);
  const [price, setPrice] = useState(-1);
  const { auth, setAuth } = useAuth();

  const instance = axios.create({
    baseURL: "http://localhost:8080",
    headers: {
      "Content-Type": "application/json",
    },
  });

  const refreshInstance = axios.create({
    baseURL: "http://localhost:8080",
    headers: {
      "Content-Type": "application/json",
    },
  });

  const directionEmojis = {
    up: '📈',
    down: '📉',
    '': '',
  };

  const chart = {
    options: {
      chart: {
        type: 'candlestick',
        height: 350
      },
      title: {
        text: 'CandleStick Chart',
        align: 'left'
      },
      xaxis: {
        type: 'datetime'
      },
      yaxis: {
        tooltip: {
          enabled: true
        }
      }
    },
  };

  const round = (number) => {
    return number ? +(number.toFixed(2)) : null;
  };

  instance.interceptors.request.use(
    async (config) => {
      const accessToken = localStorage.getItem("access_token");
      const auth = jwt_decode(accessToken);
      const expireTime = auth.exp * 1000;
      const now = + new Date();
      if (expireTime > now) {
        config.headers["Authorization"] = 'Bearer ' + accessToken;
      } else {
          const refreshToken = localStorage.getItem("refresh_token");
          const response = await refreshAccessToken(refreshToken);
          const data = response.data;
          console.log(data);
          const accessToken = data.accessToken;
          setAuth({token: accessToken,refreshToken});
          localStorage.removeItem("access_token");
          localStorage.setItem("access_token", accessToken);
          config.headers["Authorization"] = 'Bearer ' + accessToken;
          console.log("refresh is finished")
      }
      console.log("exist from interceptors")
      return config;
    },
    (error) => {
      return Promise.reject(error);
    }
  );

  const refreshAccessToken =async (token) => {
    return refreshInstance.post("/refresh", token);
  };

  const getStockData = async () => {
    return instance.get(`http://localhost:8080/api/${slug}`);
  }

  useEffect(() => {
    (async () => {
      const response = await getStockData();
      const data = response.data;
      const stockInfo = data.chart.result[0];
      console.log(stockInfo);
      setPrice(stockInfo.meta.regularMarketPrice.toFixed(2));
      setPriceTime(new Date(stockInfo.meta.regularMarketTime * 1000));
      setSymbol(stockInfo.meta.symbol);
      const quote = stockInfo.indicators.quote[0];
      const prices = stockInfo.timestamp.map((timestamp, index) => ({
        x: new Date(timestamp * 1000),
        y: [quote.open[index], quote.high[index], quote.low[index], quote.close[index]].map(round)
      }));
      setPriceInfo([{
        data: prices,
      }]);


      setStockData({ data });
    })();
  }, []);

  return (
    <div className="stock-container">
      <div>
        <h1 className="title">Last stock price</h1>
        <h1 className="price">${price} {directionEmojis['up']}</h1>
        <h2 className="symbol">{symbol}</h2>
        <h3 className="time">{priceTime && priceTime.toLocaleTimeString()}</h3>
      </div>
      <Chart options={chart.options} series={priceInfo} type="candlestick" width="100%" height={350} />
    </div>

  );

};
export default Code;
