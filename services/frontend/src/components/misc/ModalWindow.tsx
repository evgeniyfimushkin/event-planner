import React from "react";

import "./ModalWindow.css";

export default function ModalWindow({ children, onClose }: {
    children: React.JSX.Element,
    onClose: ()=>void
}) {
    return (
        <div className="modal" onClick={(e)=>e.stopPropagation()}>
            {children}
            <div className="center">
                <button onClick={onClose}>Close</button>
            </div>
        </div>
    );
}
