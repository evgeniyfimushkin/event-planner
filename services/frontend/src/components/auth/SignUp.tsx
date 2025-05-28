import "./Auth.css"

import React from "react";
import { useState, useContext } from "react";
import CryptoJS from "crypto-js";
import { useNavigate } from "react-router-dom";
import { explainRequestError } from "../../utilities/Utilities";
import { API } from "../../utilities/API";

export default function SignUp({}) {
    const [username, setUsername] = useState<string>("");
    const [email, setEmail] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const navigate = useNavigate();

    const handleSignUp = async (e) => {
        e.preventDefault();
        try {
            const passhash: string = CryptoJS.SHA256(password).toString(CryptoJS.enc.Hex);
            await API.Auth.Register({ username, email, passhash })
            //await API.Auth.Register({ username, passhash })
            alert("Зарегистрировано!")
            navigate("/signIn");
        } catch (error) {
            alert("Ошибка регистрации!\n"+explainRequestError(error));
            console.error(error);
        }
    };

    return (
        <div className="signUp">
            <h1>Регистрация</h1>
            <form onSubmit={handleSignUp} className="form">
                <input type="text" placeholder="Имя пользователя" value={username} onChange={(e) => setUsername(e.target.value)} />
                {/* <input type="email" placeholder="Почта" value={email} onChange={(e) => setEmail(e.target.value)} /> */}
                <input type="password" placeholder="Пароль" value={password} onChange={(e) => setPassword(e.target.value)} />
                <button type="submit">Зарегистрироваться</button>
            </form>
            <a href="/">Уже зарегистрированы?</a>
        </div>
    )
}
