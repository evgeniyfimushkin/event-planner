import "./ModalWindow.css";

import React from "react";

interface Button {
    name: string,
    onClick: (e?)=>any,
}

export default function ModalWindow({ children, buttons }: {
    children: React.JSX.Element,
    buttons?: Array<Button>
}) {
    return (
        <div className="modal-background">
            <div className="modal" onClick={(e)=>e.stopPropagation()}>
                {children}
                <div className="center">
                    {buttons?.map((button, index)=><button key={index} onClick={button.onClick}>
                        {button.name}
                    </button>)}
                    {/* <button onClick={onClose}>Close</button> */}
                </div>
            </div>
        </div>
    );
}
