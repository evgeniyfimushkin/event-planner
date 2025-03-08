import { createPortal } from "react-dom";
import "./Cards.css"
import { useState } from "react";
import ModalWindow from "../misc/ModalWindow";
import { explainRequestError, localDate } from "../../services/Utilities";
import axios from "axios";

export default function Card({event, subscribedInitially}) {
    const [showModal, setShowModal] = useState(false);
    const [subscribed, setSubscribed] = useState(subscribedInitially);

    const {
        id,
        name,
        description,
        category,
        participants,
        max_participants,
        image_data,
        city,
        address,
        latitude,
        longitude,
        start_time,
        end_time,
    } = event;
    const coords = (latitude && longitude) && "координаты " + latitude + " " + longitude;
    const fullAddress = [city, address, coords].filter(e=>e).join(", ");

    const subscribe = async (e) => {
        try {
            const refresh = await axios.get("/api/v1/auth/refresh");
            const res = await axios.post("/api/v1/registrations", { event_id: id });
            alert("Вы подписаны на событие!");
            setSubscribed(true);
        } catch (error) {
            alert("Не удалось подписаться!\n"+explainRequestError(error));
            console.error(error);
        }
    }
    const unsubscribe = async (e) => {
        try {
            const refresh = await axios.get("/api/v1/auth/refresh");
            const res = await axios.delete("/api/v1/registrations", { params: {event_id: id} });
            alert("Вы отписаны от события!");
            setSubscribed(false);
        } catch (error) {
            alert("Не удалось отписаться!\n"+explainRequestError(error));
            console.error(error);
        }
    }

    return (
        <>
        <div className="card" onClick={()=>setShowModal(true)}>
            <div className="line">
                {image_data && <img src={image_data || null}/>}
                <h1 className="title">{name}</h1>
                {description && <p className="description">{description}</p>}
            </div>
            {max_participants && <p className="maxParticipants">{participants}/{max_participants} участников записаны</p>}
            {fullAddress && (
                <p>Местоположение: {fullAddress}</p>
            )}
            {start_time && <p className="startTime">Начало: {localDate(new Date(start_time))}</p>}
            {end_time && <p className="endTime">Окончание: {localDate(new Date(end_time))}</p>}
            {category && <p className="category">{category}</p>}
            {/* todo refresh */}
            {subscribed && <>
                <input type="button" value="Отписаться" onClick={unsubscribe} />
            </> || <>
                <input type="button" value="Записаться" onClick={subscribe} />
            </>}
        </div>
        </>
    )
}
