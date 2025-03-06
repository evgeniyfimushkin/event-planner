import { createPortal } from "react-dom";
import "./Cards.css"
import { useState } from "react";
import ModalWindow from "../misc/ModalWindow";
import Card from "./Card";
import { localDate } from "../../services/Utilities";

export default function Minicard({event, subscribed}) {
    const [showModal, setShowModal] = useState(false);
    const {
        id,
        name,
        description,
        image_data,
        city,
        address,
        latitude,
        longitude,
        start_time,
    } = event;
    const coords = (latitude && longitude) && "координаты " + latitude + " " + longitude;
    const fullAddress = [city, address, coords].filter(e=>e).join(", ");
    return (
        <>
        <div className="card mini" onClick={()=>setShowModal(true)}>
            <div className="line">
                {image_data && <img src={image_data || null}/>}
                <h1 className="title">{name}</h1>
            </div>
            {description && <p className="description">{description}</p>}
            {fullAddress && (
                <p>Местоположение: {fullAddress}</p>
            )}
            {start_time && <p className="startTime">Начало: {localDate(new Date(start_time))}</p>}
            {showModal && createPortal(
                <ModalWindow onClose={()=>{setShowModal(false);}}>
                    <Card event={event} subscribedInitially={subscribed} />
                </ModalWindow>,
                document.body
            )}
        </div>
        </>
    )
}
