import "./Auth.css"

import { useState, useContext } from "react";
import axios from "axios";
import AuthContext from "../../services/AuthContext";

export default function Login({}) {
    const [username, setUsername] = useState("");
    const [passhash, setPasshash] = useState("");
    const { login, access_token, refresh_token } = useContext(AuthContext);

    const handleLogin = async (e) => {
        e.preventDefault();
        try {
            const res = await axios.post("http://localhost:8081/api/v1/auth/login", { username, passhash });

            login({
                "access_token": res.data.access_token,
                "refresh_token": res.data.refresh_token,
            });
            alert("Connected!");
            console.log(res.headers);
        } catch (error) {
            alert("Connection error!\n"+e.message);
        }
    };

    return (
        <div>
            <h1>Log in</h1>
            <form onSubmit={handleLogin} className="form">
                <input type="text" placeholder="Username" value={username} onChange={(e) => setUsername(e.target.value)} />
                <input type="password" placeholder="Password" value={passhash} onChange={(e) => setPasshash(e.target.value)} />
                <button type="submit">Log in</button>
            </form>
        </div>
    )
}
