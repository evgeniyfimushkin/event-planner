import "./ModalWindow.css";

export default function ModalWindow({ children, onClose }) {
    return (
        <div className="modal" onClick={(e)=>e.stopPropagation()}>
            {children}
            <div className="center">
                <button onClick={onClose}>Close</button>
            </div>
        </div>
    );
}
