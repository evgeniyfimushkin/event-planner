import React from "react";
import { useContext } from "react";
import { Navigate } from "react-router-dom";
import AuthContext from "./AuthContext";

export default function PrivateRoute({ children } : {
    children: React.JSX.Element | React.JSX.Element[]
}) {
    const { signedIn } = useContext(AuthContext);
    console.log(signedIn);
    return signedIn ? children : <Navigate to="/signIn" />;
};
