import "./Auth.css"

import { useState, useContext } from "react";
import axios from "axios";
import CryptoJS from "crypto-js";
import { useNavigate } from "react-router-dom";

export default function Register({}) {
    const [username, setUsername] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const navigate = useNavigate();

    const handleRegister = async (e) => {
        e.preventDefault();
        try {
            const passhash = CryptoJS.SHA256(password).toString(CryptoJS.enc.Hex);
            const res = await axios.post("http://localhost/api/v1/auth/register", { username, email, passhash });
            alert("Registered!")
            navigate("/login");
        } catch (error) {
            alert("Cannot register!\n"+e);
        }
    };

    return (
        <div>
            <h1>Register</h1>
            <form onSubmit={handleRegister} className="form">
                <input type="text" placeholder="Username" value={username} onChange={(e) => setUsername(e.target.value)} />
                <input type="email" placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} />
                <input type="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} />
                <button type="submit">Register</button>
            </form>
        </div>
    )
}
