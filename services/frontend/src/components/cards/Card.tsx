import React, { useContext } from "react";

import { createPortal } from "react-dom";
import "./Cards.css"
import { useState } from "react";
import ModalWindow from "../misc/ModalWindow";
import { explainRequestError, localDate } from "../../utilities/Utilities";
import axios from "axios";
import { useEffect } from "react";
import { authCall } from "../../utilities/Utilities";
import { Event } from "../../utilities/Types";
import { useNavigate } from "react-router-dom";
import AuthContext from "../../services/AuthContext";

export default function Card({event}: {
    event: Event
}) {
    const [showModal, setShowModal] = useState<boolean>(false);
    const [subscribed, setSubscribed] = useState<boolean>(false);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);
    const { logout } = useContext(AuthContext);
    const navigate = useNavigate();

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
    const coords: string | false = (["latitude", "longitude"].every(e=>e in event)) && "координаты " + latitude + " " + longitude;
    const fullAddress: string = [city, address, coords].filter(e=>e).join(", ");

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
            const res = await axios.delete("/api/v1/registrations", { data: {event_id: id} });
            alert("Вы отписаны от события!");
            setSubscribed(false);
        } catch (error) {
            alert("Не удалось отписаться!\n"+explainRequestError(error));
            console.error(error);
        }
    }

    const fetchSubscribed = async () => {
        try {
            await authCall(async () => { // success
                setLoading(true);
                const responseSubscriptions = await axios.get("/api/v1/registrations/my");
                setSubscribed(responseSubscriptions.data.some(s=>s.event_id === id));
            }, (err) => { // unauthorized
                logout();
                navigate("/login");
            });
        } catch (err) {
            setError("Не удалось загрузить статус записи!\n"+explainRequestError(err));
            console.error(err);
        } finally {
            setLoading(false);
        }
    }

    useEffect(() => {
        fetchSubscribed();
    }, []);

    // if (loading) return <p>Загрузка...</p>;
    // if (error) return <p>{error}</p>;

    return (
        <>
        <div className="card" onClick={()=>setShowModal(true)}>
            <div className="line">
                {image_data && <img src={image_data}/>}
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
            {!loading && !error && (subscribed && <>
                <input type="button" value="Отписаться" onClick={unsubscribe} />
            </> || <>
                <input type="button" value="Записаться" onClick={subscribe} />
            </>)}
        </div>
        </>
    )
}
