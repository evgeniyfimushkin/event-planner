import { createContext, useState, useEffect } from "react";
import { useCookies } from "react-cookie";

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
    // const [token, setToken] = useState(localStorage.getItem("token") || null);
    const [cookies, setCookie, removeCookie] = useCookies(["access_token", "refresh_token"]);
    // const [setAccessToken] = useState(cookies?.access_token || null);
    // const [setRefreshToken] = useState(cookies?.refresh_token || null);
    const [loggedIn, setLoggedIn] = useState(false);

    const login = () => {
        // setToken(newToken);
        // localStorage.setItem("token", newToken);
        // todo set expiration

        // setAccessToken(newTokens.access_token);
        // setRefreshToken(newTokens.refresh_token);
        // setCookie("access_token", newTokens.access_token, {path: "/"});
        // setCookie("refresh_token", newTokens.refresh_token, {path: "/"});
        setLoggedIn(true);
    };

    const logout = () => {
        // setToken(null);
        // localStorage.removeItem("token");

        // setAccessToken(null);
        // setRefreshToken(null);
        removeCookie("access_token");
        removeCookie("refresh_token");
        setLoggedIn(false);
    };

    return (
        <AuthContext.Provider value={{ loggedIn, login, logout }}>
            {children}
        </AuthContext.Provider>
    );
};

export default AuthContext;
