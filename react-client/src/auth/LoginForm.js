import React, { useEffect, useState } from 'react'
import {useNavigate,useLocation} from "react-router-dom";
import axios from "axios";
import useAuth from '../routing/useAuth'
import "./LoginForm.css"
import 'bootstrap/dist/css/bootstrap.css';

const Login = () =>{

  const errors = {
    credentials:"The login or password is incorrect",
    internal:"Internal server error",
  };
  const [fieldError,setFieldError] = useState("");
  const navigate = useNavigate();
  const {auth,setAuth} = useAuth();
  const location = useLocation();
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

  const getToken = async (userData) =>{
    return instance
    .post("/api/auth/login", userData)
    .catch(error => {   
      if (error.response.status === 400){
        setFieldError(errors.credentials)
        console.log("field = " + fieldError)
      }else{
        setFieldError(errors.internal)
      }
    });
  }

  useEffect(() => {
    if (auth.token) {
      localStorage.setItem("access_token", auth.token); 
      navigate("/home");
    }
  }, [auth]);

  useEffect(() => {
    if (location.state){
      console.log(location.state.previousUrl)
    }
  }, []);

  const onLogin = (e) =>{
    e.preventDefault();
    location.state =null
    const userData = {
      username: login,
      password: password,
    };

    (async() => {
      const response = await getToken(userData);
      const data =  response.data;
      const token = data.access_token;
      setAuth({token});
    })();
    
};


  return (
    <div className="Auth-form-container">
      <form className="Auth-form" onSubmit={onLogin}>
        <div className="Auth-form-content">
          <h3 className="Auth-form-title">Sign In</h3>
          <div className="form-group mt-3">
                {location.state && (
                          <div className="alert alert-success" role="alert">
                          {"you have been successfully registered"}
                          </div>
                )}
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
                {fieldError && (
                            <div className="alert alert-danger" role="alert">
                            {fieldError}
                            </div>
                )}
          </div>
        </div>
      </form>
    </div>
  )
};
export default Login;

