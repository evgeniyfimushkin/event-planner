import express from "express";
import jwt from "jsonwebtoken";
import http from "http";

const app = express();
const server = http.createServer(app);

app.use(express.json());

const SECRET = "thesecret";

app.post("/api/v1/auth/login", (req, res) => {
    const { email, password } = req.body;
    if (email === "ivan" && password === "1234") {
        const token = jwt.sign({ email }, SECRET, { expiresIn: "1h" });
        return res.json({ token });
    }
    res.status(401).json({ error: "Неверные данные!" });
});

server.listen(5000, () => console.log("Server on port 5000"));
