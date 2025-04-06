import React from "react";

import { createPortal } from "react-dom";
import "./Cards.css"
import { useState } from "react";
import ModalWindow from "../misc/ModalWindow";
import Card from "./Card";
import { localDate } from "../../utilities/Utilities";
import { Event } from "../../utilities/Types";

export default function Minicard({event}: {
    event: Event
}) {
    const [showModal, setShowModal] = useState<boolean>(false);
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
    const coords: string | false = (["latitude", "longitude"].every(e=>e in event)) && "координаты " + latitude + " " + longitude;
    const fullAddress: string = [city, address, coords].filter(e=>e).join(", ");
    return (
        <>
        <div className="card mini" onClick={()=>setShowModal(true)}>
            <div className="line">
                {image_data && <img src={image_data}/>}
                <h1 className="title">{name}</h1>
                {description && <p className="description">{description}</p>}
            </div>
            {fullAddress && (
                <p>Местоположение: {fullAddress}</p>
            )}
            {start_time && <p className="startTime">Начало: {localDate(new Date(start_time))}</p>}
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
