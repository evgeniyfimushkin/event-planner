import "./Cards.css"

import React, { useContext } from "react";
import { useState } from "react";
import { explainRequestError, localeDateString } from "../../utilities/Utilities";
import { useEffect } from "react";
import { authCall } from "../../utilities/Utilities";
import { Event, Review, UserData } from "../../utilities/Types";
import { useNavigate } from "react-router-dom";
import AuthContext from "../../services/AuthContext";
import { API } from "../../utilities/API";
import Reviews from "../reviews/Reviews";

export default function Card({event}: {
    event: Event
}) {
    const [subscribed, setSubscribed] = useState<boolean>(false);
    const [loadingRegistration, setLoadingRegistration] = useState<boolean>(true);
    const [errorRegistration, setErrorRegistration] = useState<string | null>(null);

    const [user, setUser] = useState<UserData | null>(null);
    const [loadingUser, setLoadingUser] = useState<boolean>(true);
    const [errorUser, setErrorUser] = useState<string | null>(null);

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
        created_by,
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
                setLoadingRegistration(true);
                const responseSubscriptions = await API.Registrations.GetMy();
                setSubscribed(responseSubscriptions.data.some(s=>s.event_id === id));
            }, (err) => { // unauthorized
                signOut();
                navigate("/signIn");
            });
        } catch (err) {
            setErrorRegistration("Не удалось загрузить статус записи!\n"+explainRequestError(err));
            console.error(err);
        } finally {
            setLoadingRegistration(false);
        }
    }
    useEffect(() => {
        fetchSubscribed();
    }, []);

    const fetchUser = async () => {
        try {
            await authCall(async () => { // success
                setLoadingUser(true);
                const responseData = await API.Users.GetData(created_by!);
                setUser(responseData.data);
            }, (err) => { // unauthorized
                // signOut();
                // navigate("/signIn");
                throw new Error("Unauthorized"); // is this applicable?
            });
        } catch (err) {
            setErrorUser("Не удалось загрузить данные пользователя!\n"+explainRequestError(err));
            console.error(err);
        } finally {
            setLoadingUser(false);
        }
    }
    useEffect(() => {
        fetchUser();
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
            {start_time && <p className="startTime">Начало: {localeDateString(new Date(start_time))}</p>}
            {end_time && <p className="endTime">Окончание: {localeDateString(new Date(end_time))}</p>}
            {category && <p className="category">Категория: {category}</p>}
            <p className="user">{
                (loadingUser && <></>) ||
                (errorUser && <>Пользователь <span className="username">{created_by}</span></> ) ||
                <>
                    {user!.picture && <img src={user!.picture}/>}
                    <span className="username">{user!.username}</span>
                </>
            }</p>
            {/* todo refresh */}
            {!loadingRegistration && !errorRegistration && (subscribed && <>
                <input type="button" value="Отписаться" onClick={unsubscribe} />
            </> || <>
                <input type="button" value="Записаться" onClick={subscribe} />
            </>)}
            <Reviews event_id={id!} />
        </div>
        </>
    )
}
