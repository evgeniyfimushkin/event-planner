import "./Auth.css"

import { useState, useContext } from "react";
import axios from "axios";

export default function Register({}) {
    const [username, setUsername] = useState("");
    const [email, setEmail] = useState("");
    const [passhash, setPasshash] = useState("");

    const handleRegister = async (e) => {
        e.preventDefault();
        try {
            const res = await axios.post("http://localhost:5000/api/v1/auth/register", { username, email, passhash });
            // todo process response
            alert("Registered!")
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
                <input type="password" placeholder="Password" value={passhash} onChange={(e) => setPasshash(e.target.value)} />
                <button type="submit">Register</button>
            </form>
        </div>
    )
}
