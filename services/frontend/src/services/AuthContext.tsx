import React from "react";
import { createContext, useState, useEffect } from "react";
import { UserCredentials } from "../utilities/Types";



interface AuthContextType {
    signedIn: boolean,
    username: string,
    signIn: (params: {username: string})=>void,
    signOut: ()=>void,
};

const AuthContext = createContext<AuthContextType>({
    signedIn: false,
    username: "",
    signIn: ()=>console.error("signIn(): Некорректное использование AuthContext"),
    signOut: ()=>console.error("signOut(): Некорректное использование AuthContext"),
});

export const AuthProvider = ({ children }) => {
    const [signedIn, setSignedIn] = useState(!!localStorage.getItem("signedIn") || false);
    const [username, setUsername] = useState(localStorage.getItem("username") ?? "");

    const signIn = ({username}) => {
        setSignedIn(true);
        setUsername(username);
        localStorage.setItem("signedIn", ""+true);
        localStorage.setItem("username", username);
    };
    const signOut = () => {
        setUsername("");
        //setSignedIn(false);
        localStorage.removeItem("signedIn");
        localStorage.removeItem("username");
    };

    return (
        <AuthContext.Provider value={{ username, signedIn, signIn, signOut }}>
            {children}
        </AuthContext.Provider>
    );
};

export default AuthContext;
