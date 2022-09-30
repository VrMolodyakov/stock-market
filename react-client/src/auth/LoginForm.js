import React, { Component, useEffect, useState } from 'react'
import {useNavigate} from "react-router-dom";
import axios from "axios";
import useAuth from '../routing/useAuth'
import 'bootstrap/dist/css/bootstrap.css';

const Login = () =>{

  const navigate = useNavigate();
  const {auth,setAuth} = useAuth();

  const [authResponce,setAuthResponce] = useState(null);
  const [login,setLogin] = useState("");
  const onChangeLogin = (e) => {
        e.preventDefault();
        const username = e.target.value;
        setLogin(username);
  };

  const [password,setPassword] = useState("");
  const onChangePassword = (e) => {
        const password = e.target.value;
        setPassword(password);
  }; 

  const instance = axios.create({
    baseURL: "http://localhost:8080",
    withCredentials: true,
    headers: {
      "Content-Type": "application/json",
    },
  });

  const refreshAccessToken = () =>{
    console.log("refresh call");
  }

  instance.interceptors.response.use((response) => {
    return response
  }, async function (error) {
    const originalRequest = error.config;
    if (error.response.status === 403 && !originalRequest._retry) {
      originalRequest._retry = true;
      const access_token = await refreshAccessToken();            
      axios.defaults.headers.common['Authorization'] = 'Bearer ' + access_token;
      return instance(originalRequest);
    }
    return Promise.reject(error);
  });


  const getToken = async (userData) =>{
    return instance.post("/api/auth/login", userData);
  }

  useEffect(() => {
    if (auth.token) {
      console.log(auth);
      localStorage.setItem("access_token", auth.token); 
      navigate("/home");
    }
  }, [auth]);

  const onLogin = (e) =>{
    e.preventDefault();

    const userData = {
      username: login,
      password: password,
    };

    (async() => {
      const response = await getToken(userData);
      const data =  response.data;
      const token = data.access_token;
      console.log(token)
      const refreshToken = data.refreshToken;
      setAuth({token});
      console.log(auth);
    })();
    
};


  return (
    <div className="Auth-form-container">
      <form className="Auth-form" onSubmit={onLogin}>
        <div className="Auth-form-content">
          <h3 className="Auth-form-title">Sign In</h3>
          <div className="form-group mt-3">
            <label>Login</label>
            <input
              type="login"
              className="form-control mt-1"
              placeholder="Enter login"
              value={login}
              onChange={onChangeLogin}
            />
          </div>
          <div className="form-group mt-3">
            <label>Password</label>
            <input
              type="password"
              className="form-control mt-1"
              placeholder="Enter password"
              value={password}
              onChange={onChangePassword}
            />
          </div>
          <div className="d-grid gap-2 mt-3">
            <button type="submit" className="btn btn-primary">
              Submit
            </button>
          </div>
        </div>
      </form>
    </div>
  )
};
export default Login;
