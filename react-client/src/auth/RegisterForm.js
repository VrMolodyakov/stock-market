import React, { Component, useEffect, useState } from 'react'
import 'bootstrap/dist/css/bootstrap.css';
import axios from "axios";

const Register = () => {

    const onRegister = (e) =>{
        e.preventDefault();

        const userData = {
            login: login,
            password: password,
            email:email
        };
        axios.post("http://localhost:8080/registration", userData)
             .then((response) => {
                console.log(response.status);
                console.log(response);
        });
    };


    const [login,setLogin] = useState("");
    const onChangeLogin = (e) => {
        const username = e.target.value;
        setLogin(username);
    };

    const [password,setPassword] = useState("");
    const onChangePassword = (e) => {
        const password = e.target.value;
        setPassword(password);
    }; 

    const [email,setEmail] = useState("");
    const onChangeEmail = (e) => {
        const email = e.target.value;
        setEmail(email);
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
                </div>
                <div className="form-group mt-3">
                <label>Email address</label>
                <input
                    type="email"
                    className="form-control mt-1"
                    placeholder="Horus@mail.com"
                    value={email}
                    onChange={onChangeEmail}
                />
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
                </div>
                <div className="d-grid gap-2 mt-3">
                <button type="submit" className="btn btn-primary">
                    Submit
                </button>
                </div>
                <p className="text-center mt-2">
                Forgot <a href="#">password?</a>
                </p>
            </div>
            </form>
        </div>
    )
};
export default Register;
