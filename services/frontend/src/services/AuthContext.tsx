import React from "react";

import { createContext, useState, useEffect } from "react";
import { useCookies } from "react-cookie";

interface AuthContextType {
    signedIn: boolean,
    signIn: ()=>void,
    signOut: ()=>void,
};

const AuthContext = createContext<AuthContextType>({
    signedIn: false,
    signIn: ()=>console.error("signIn(): Некорректное использование AuthContext"),
    signOut: ()=>console.error("signOut(): Некорректное использование AuthContext"),
});

export const AuthProvider = ({ children }) => {
    // const [token, setToken] = useState(localStorage.getItem("token") || null);
    const [cookies, setCookie, removeCookie] = useCookies(["access_token", "refresh_token"]);
    // const [setAccessToken] = useState(cookies?.access_token || null);
    // const [setRefreshToken] = useState(cookies?.refresh_token || null);
    const [signedIn, setSignedIn] = useState(!!localStorage.getItem("signedIn") || false);

    const signIn = () => {
        // setToken(newToken);
        // localStorage.setItem("token", newToken);
        // todo set expiration

        // setAccessToken(newTokens.access_token);
        // setRefreshToken(newTokens.refresh_token);
        // setCookie("access_token", newTokens.access_token, {path: "/"});
        // setCookie("refresh_token", newTokens.refresh_token, {path: "/"});
        setSignedIn(true);
        localStorage.setItem("signedIn", ""+true);
    };

    const signOut = () => {
        // setToken(null);
        // localStorage.removeItem("token");

        // setAccessToken(null);
        // setRefreshToken(null);
        removeCookie("access_token");
        removeCookie("refresh_token");
        setSignedIn(false);
        localStorage.removeItem("signedIn");
    };

    return (
        <AuthContext.Provider value={{ signedIn, signIn, signOut }}>
            {children}
        </AuthContext.Provider>
    );
};

export default AuthContext;
