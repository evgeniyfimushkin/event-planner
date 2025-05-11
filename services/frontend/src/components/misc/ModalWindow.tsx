import "./ModalWindow.css";

import React from "react";

export default function ModalWindow({ children, onClose }: {
    children: React.JSX.Element,
    onClose: ()=>void
}) {
    return (
        <div className="modal-background">
            <div className="modal" onClick={(e)=>e.stopPropagation()}>
                {children}
                <div className="center">
                    <button onClick={onClose}>Close</button>
                </div>
            </div>
        </div>
    );
}
