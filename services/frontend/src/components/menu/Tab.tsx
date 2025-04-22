import "./Menu.css";

import React from "react";

export default function Tab({title, target, onClick, isOpen}: {
    title: string,
    target?: string,
    onClick?: ()=>void,
    isOpen?: boolean,
}) {
    return (
        <div className={"tab " + (isOpen && "open" || "")} onClick={onClick || (()=>{})}>
            {target
            && <a href={target}>{title}</a>
            || <span>{title}</span>}
        </div>
    )
}
