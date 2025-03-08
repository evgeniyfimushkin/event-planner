import "./Auth.css"

import { useState, useContext } from "react";
import axios from "axios";
import AuthContext from "../../services/AuthContext";
import { useNavigate } from "react-router-dom";
import CryptoJS from "crypto-js";
import { explainRequestError } from "../../services/Utilities";

export default function Login({}) {
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const { login } = useContext(AuthContext);
    const navigate = useNavigate();

    const handleLogin = async (e) => {
        e.preventDefault();
        try {
            const passhash = CryptoJS.SHA256(password).toString(CryptoJS.enc.Hex);
            const res = await axios.post("/api/v1/auth/login", { username, passhash });

            // const fr = await fetch("http://localhost/api/v1/auth/login", {
            //     method: "POST",
            //     headers: {
            //         'Content-Type': 'application/json',
            //     },
            //     body: JSON.stringify({ username, passhash }),
            // });

            // login({
            //     "access_token": res.data.access_token,
            //     "refresh_token": res.data.refresh_token,
            // });
            login();
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
        <div>
            <h1>Вход</h1>
            <form onSubmit={handleLogin} className="form">
                <input type="text" placeholder="Имя пользователя" value={username} onChange={(e) => setUsername(e.target.value)} />
                <input type="password" placeholder="Пароль" value={password} onChange={(e) => setPassword(e.target.value)} />
                <button type="submit">Log in</button>
            </form>
            <a href="/register">Впервые здесь?</a>
        </div>
    )
}
