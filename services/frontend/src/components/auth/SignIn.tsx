import React from "react";

import "./Auth.css"

import { useState, useContext } from "react";
import axios from "axios";
import AuthContext from "../../services/AuthContext";
import { useNavigate } from "react-router-dom";
import CryptoJS from "crypto-js";
import { explainRequestError } from "../../utilities/Utilities";

export default function SignIn({}) {
    const [username, setUsername] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const { signIn } = useContext(AuthContext);
    const navigate = useNavigate();

    const handleSignIn = async (e) => {
        e.preventDefault();
        try {
            const passhash: string = CryptoJS.SHA256(password).toString(CryptoJS.enc.Hex);
            const res = await axios.post("/api/v1/auth/login", { username, passhash });

            // const fr = await fetch("http://localhost/api/v1/auth/login", {
            //     method: "POST",
            //     headers: {
            //         'Content-Type': 'application/json',
            //     },
            //     body: JSON.stringify({ username, passhash }),
            // });

            // signIn({
            //     "access_token": res.data.access_token,
            //     "refresh_token": res.data.refresh_token,
            // });
            signIn();
            alert("Подключено!");
            const refresh = await axios.get("/api/v1/auth/refresh");
            navigate("/");
            // console.log(res.headers);
        } catch (error) {
            alert("Ошибка подключения!\n"+explainRequestError(error));
            console.error(error);
        }
    };

    return (
        <div className="signIn">
            <h1>Вход</h1>
            <form onSubmit={handleSignIn} className="form">
                <input type="text" placeholder="Имя пользователя" value={username} onChange={(e) => setUsername(e.target.value)} />
                <input type="password" placeholder="Пароль" value={password} onChange={(e) => setPassword(e.target.value)} />
                <button type="submit">Войти</button>
            </form>
            <a href="/signUp">Впервые здесь?</a>
        </div>
    )
}
