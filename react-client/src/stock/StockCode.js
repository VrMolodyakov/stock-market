import { useParams,useNavigate } from "react-router-dom";
import useAuth from '../routing/useAuth';
import { useState, useEffect, useMemo } from 'react'
import Logout from "../auth/Logout"
import Chart from 'react-apexcharts';
import axios from "axios";
import jwt_decode from 'jwt-decode'
import "./StockCode.css"

const Code = (props) => {
  const { slug } = useParams();
  const navigate = useNavigate();
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
    withCredentials: true,
    headers: {
      "Content-Type": "application/json",
    },
  });

  const refreshInstance = axios.create({
    baseURL: "http://localhost:8080",
    withCredentials: true,
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
          const response = await refreshAccessToken();
          const data = response.data;
          const accessToken = data.access_token;
          setAuth({token: accessToken});
          localStorage.removeItem("access_token");
          localStorage.setItem("access_token", accessToken);
          config.headers["Authorization"] = 'Bearer ' + accessToken;
      }
      console.log("exist from interceptors")
      return config;
    },
    (error) => {
      console.log(error)
      console.log("token is expired")
      }
  );

  const refreshAccessToken =async () => {
    return refreshInstance.get("/api/auth/refresh");
  };

  const getStockData = async () => {
    return instance.get(`/api/stock/symbols/${slug}`);
  }

  useEffect(() => {
    (async () => {
      const response = await getStockData();
      console.log(response)
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
    })().catch( 
      (error) =>{
        console.log(error)
        Logout(auth, setAuth)
        navigate("/auth");
      } 
    );
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
