import React from "react";
import { createContext, useState, useEffect } from "react";

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
    const [signedIn, setSignedIn] = useState(!!localStorage.getItem("signedIn") || false);

    const signIn = () => {
        setSignedIn(true);
        localStorage.setItem("signedIn", ""+true);
    };
    const signOut = () => {
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
