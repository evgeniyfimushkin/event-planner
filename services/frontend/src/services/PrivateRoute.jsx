import { useContext } from "react";
import { Navigate } from "react-router-dom";
import AuthContext from "./AuthContext";

const PrivateRoute = ({ children }) => {
    const { accessToken } = useContext(AuthContext);
    return accessToken ? children : <Navigate to="/login" />;
};

export default PrivateRoute;
