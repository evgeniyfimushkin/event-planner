import "./Event.css"

import React from "react";
import { useState, useContext } from "react";
import { explainRequestError, toLocalDate } from "../../utilities/Utilities";
import { Event } from "../../utilities/Types";
import { API } from "../../utilities/API";

export default function EditEvent({event, callback}: {
    event: Event,
    callback?: (e?) => any
}) {
    const [name, setName] = useState(event.name);
    const [description, setDescription] = useState(event.description);
    const [category, setCategory] = useState(event.category);
    const [maxParticipants, setMaxParticipants] = useState(event.max_participants);
    const [imageData, setImageData] = useState(event.image_data ?? "");
    const [city, setCity] = useState(event.city ?? "");
    const [address, setAddress] = useState(event.address ?? "");
    const [latitude, setLatitude] = useState(event.latitude ?? 0);
    const [longitude, setLongitude] = useState(event.longitude ?? 0);
    const [startTime, setStartTime] = useState(toLocalDate(new Date(event.start_time)));
    const [endTime, setEndTime] = useState(toLocalDate(new Date(event.end_time)));

    const updateEvent = async (e) => {
        e.preventDefault();
        try {
            await API.Auth.Refresh();
            const pack: Event = {
                id: event.id,
                name,
                description,
                category,
                max_participants: +maxParticipants,
                image_data: imageData,
                city,
                address,
                latitude: +latitude,
                longitude: +longitude,
                start_time: new Date(startTime).toISOString(),
                end_time: new Date(endTime).toISOString(),
            };
            const res = await API.Events.Update(pack);
            alert("Мероприятие обновлено!");
        } catch (error) {
            alert("Ошибка обновления мероприятия!\n" + explainRequestError(error));
            console.error(error)
        }
    };

    const handleImageData = (e) => {
        const file = e.target.files[0];
        if (!file) return;
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = () => {
            setImageData(reader.result as string);
        };
        reader.onerror = (pe) => {
            alert("Невозможно загрузить изображение!\n"+reader.error)
            console.error(reader.error);
        };
    };

    return (
        <div className="create-event">
            <h1>Изменить мероприятие</h1>
            <form onSubmit={updateEvent} className="form">
                    <label htmlFor="name">Название:</label>
                    <input id="name" type="text" value={name} onChange={e=>setName(e.target.value)} required />
                    <label htmlFor="desctiption">Описание:</label>
                    <textarea id="description" value={description} onChange={e=>setDescription(e.target.value)} required />
                    <label htmlFor="category">Категория:</label>
                    <input id="category" type="text" value={category} onChange={e=>setCategory(e.target.value)} required />
                    <label htmlFor="maxParticipants">Участники:</label>
                    <input id="maxParticipants" type="number" min="1" value={maxParticipants} onChange={e=>setMaxParticipants(+e.target.value)} required />
                    <label htmlFor="imageData">Иконка:</label>
                    <input id="imageData" type="file" accept="image/*" onChange={e=>handleImageData(e)} />
                    <label htmlFor="city">Город:</label>
                    <input id="city" type="text" value={city} onChange={e=>setCity(e.target.value)} />
                    <label htmlFor="address">Адрес:</label>
                    <input id="address" type="text" value={address} onChange={e=>setAddress(e.target.value)} />
                    <label htmlFor="latitude">Широта:</label>
                    <input id="latitude" type="number" min="-90" max="90" value={latitude} onChange={e=>setLatitude(+e.target.value)} />
                    <label htmlFor="longitude">Долгота:</label>
                    <input id="longitude" type="number" value={longitude} min="-180" max="180" onChange={e=>setLongitude(+e.target.value)} />
                    <label htmlFor="startTime">Время начала:</label>
                    <input id="startTime" type="datetime-local" value={startTime} onChange={e=>setStartTime(e.target.value)} required />
                    <label htmlFor="endTime">Время окончания:</label>
                    <input id="endTime" type="datetime-local" value={endTime} onChange={e=>setEndTime(e.target.value)} />
                <button type="submit">Применить</button>
            </form>
        </div>
    )
}
