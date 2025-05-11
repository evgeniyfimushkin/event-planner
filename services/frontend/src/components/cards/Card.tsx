import "./Cards.css"

import React, { useContext } from "react";
import { useState } from "react";
import { explainRequestError, localDate } from "../../utilities/Utilities";
import { useEffect } from "react";
import { authCall } from "../../utilities/Utilities";
import { Event } from "../../utilities/Types";
import { useNavigate } from "react-router-dom";
import AuthContext from "../../services/AuthContext";
import { API } from "../../utilities/API";

export default function Card({event}: {
    event: Event
}) {
    const [subscribed, setSubscribed] = useState<boolean>(false);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);
    const { signOut } = useContext(AuthContext);
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
            await API.Auth.Refresh();
            await API.Registrations.Create(id!);
            alert("Вы подписаны на событие!");
            setSubscribed(true);
        } catch (error) {
            alert("Не удалось подписаться!\n"+explainRequestError(error));
            console.error(error);
        }
    }
    const unsubscribe = async (e) => {
        try {
            await API.Auth.Refresh();
            await API.Registrations.Delete(id!);
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
                const responseSubscriptions = await API.Registrations.GetMy();
                setSubscribed(responseSubscriptions.data.some(s=>s.event_id === id));
            }, (err) => { // unauthorized
                signOut();
                navigate("/signIn");
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

    return (
        <>
        <div className="card">
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
