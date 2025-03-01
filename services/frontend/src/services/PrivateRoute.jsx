import { useContext } from "react";
import { Navigate } from "react-router-dom";
import AuthContext from "./AuthContext";

const PrivateRoute = ({ children }) => {
    const { loggedIn } = useContext(AuthContext);
    console.log(loggedIn);
    return loggedIn ? children : <Navigate to="/login" />;
};

export default PrivateRoute;
