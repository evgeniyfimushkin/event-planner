import "./FloatingButton.css";

import React from "react";
import { useNavigate } from "react-router-dom";

export default function FloatingButton({text, onClick, target}: {
    text: string,
    onClick?: ()=>void,
    target?: string
}) {
    const navigate = useNavigate();
    const redirect = () => {
        navigate(target as string);
    }

    return (
        <input className="floatingButton" type="button" value={text} onClick={()=>{
            target && redirect();
            onClick && onClick();
        }} />
    );
}
