import React, { useState,useRef, useEffect } from 'react'
import {useNavigate,useLocation} from "react-router-dom";
import { validLogin, validPassword } from './Validate';
import 'bootstrap/dist/css/bootstrap.css';
import axios from "axios";

const Register = () => {
    const navigate = useNavigate();
    const [errorMessages, setErrorMessages] = useState({});
    const loginErrorRef = useRef(false);
    const pwdErrorRef = useRef(false);
    const [loginError,setLoginError] = useState(false);
    const [pwdError,setPwdError] = useState(false);
    const [login,setLogin] = useState("");
    const [password,setPassword] = useState("");
    const [fieldError,setFieldError] = useState("");
    const location = useLocation();

    const errors = {
        login:"login must be at least 3 characters and contain : 0-9,a-z,A-Z",
        password:"password must be at least 3 characters and contain at least 1 symbol and 1 digit",
    };

    const validate = () =>{
        loginErrorRef.current = false;
        pwdErrorRef.current = false;

        if (!validLogin.test(login)) {
            loginErrorRef.current = true;
        }  
        if (!validPassword.test(password)) {
            pwdErrorRef.current = true;
        }
    };

    const onRegister = (e) =>{
        e.preventDefault();
        setFieldError("");

        validate();

        if (!loginErrorRef.current && !pwdErrorRef.current){
            const userData = {
                username: login,
                password: password,
            };
            axios.post("http://localhost:8080/api/auth/register", userData)
                .then(() => {
                    navigate("/auth",{ state: { previousUrl: location.pathname } });
                })
                .catch(error => {   
                    console.log(error)
                    setFieldError(error.response.data);
                })
        }else{
                setLoginError(loginErrorRef.current);
                setPwdError(pwdErrorRef.current);
        } 
    };

   
    const onChangeLogin = (e) => {
        const username = e.target.value;
        setLogin(username);
    };

    const onChangePassword = (e) => {
        const password = e.target.value;
        setPassword(password);
    }; 

    return (
        <div className="Auth-form-container">
            <form className="Auth-form" onSubmit={onRegister}>
            <div className="Auth-form-content">
                <h3 className="Auth-form-title">Sign up</h3>
                <div className="form-group mt-3">
                    <label>Login</label>
                    <input
                        type="login"
                        className="form-control mt-1"
                        placeholder="e.g Horus"
                        value={login}
                        onChange={onChangeLogin}
                        
                    />
                    {loginError && (
                        <div className="alert alert-danger" role="alert">
                        {errors.login}
                    </div>
                    )}
                </div>
                <div className="form-group mt-3">
                    <label>Password</label>
                    <input
                        type="password"
                        className="form-control mt-1"
                        placeholder="Password"
                        value={password}
                        onChange={onChangePassword}
                    />
                     {pwdError && (
                        <div className="alert alert-danger" role="alert">
                            {errors.password}
                        </div>
                        
                     )}
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
export default Register;
