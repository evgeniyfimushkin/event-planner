import { createPortal } from "react-dom";
import "./Cards.css"
import { useState } from "react";
import ModalWindow from "../misc/ModalWindow";
import Card from "./Card";

export default function Minicard({event}) {
    const [showModal, setShowModal] = useState(false);
    return (
        <>
        <div className="card" onClick={()=>setShowModal(true)}>
            <div className="line">
                {event.image_data && <img src={event.image_data || null}/>}
                <h1 className="title">{event.name}</h1>
                {event.description && <p className="description">{event.description}</p>}
            </div>
            {showModal && createPortal(
                <ModalWindow onClose={()=>{setShowModal(false);}}>
                    <Card event={event} />
                </ModalWindow>,
                document.body
            )}
        </div>
        </>
    )
}
