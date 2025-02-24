import "./Auth.css"

import { useState, useContext } from "react";
import axios from "axios";
import AuthContext from "../../services/AuthContext";

export default function Login({}) {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const { login } = useContext(AuthContext);

    const handleLogin = async (e) => {
        e.preventDefault();
        try {
            const res = await axios.post("http://localhost:5000/api/v1/auth/login", { email, password });
            login(res.data.token);
            alert("Connected!")
        } catch (error) {
            alert("Connection error!\n"+e);
        }
    };

    return (
        <div>
            <h1>Log in</h1>
            <form onSubmit={handleLogin} className="form">
                <input type="text" placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} />
                <input type="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} />
                <button type="submit">Log in</button>
            </form>
        </div>
    )
}
