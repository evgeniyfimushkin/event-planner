import "./Auth.css"

import React from "react";
import { useState, useContext } from "react";
import AuthContext from "../../services/AuthContext";
import { useNavigate } from "react-router-dom";
import CryptoJS from "crypto-js";
import { explainRequestError } from "../../utilities/Utilities";
import { API } from "../../utilities/API";

export default function SignIn({}) {
    const [username, setUsername] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const { signIn } = useContext(AuthContext);
    const navigate = useNavigate();

    const handleSignIn = async (e) => {
        e.preventDefault();
        try {
            const passhash: string = CryptoJS.SHA256(password).toString(CryptoJS.enc.Hex);
            await API.Auth.Login({username, passhash});
            signIn();
            await API.Auth.Refresh();
            navigate("/");
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
