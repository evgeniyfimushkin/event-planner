import "./Event.css"

import { useState, useContext } from "react";
import axios from "axios";
import { explainRequestError } from "../../services/Utilities";

export default function CreateEvent({}) {
    const [name, setName] = useState("Без названия");
    const [description, setDescription] = useState("Нет описания");
    const [category, setCategory] = useState("");
    const [maxParticipants, setMaxParticipants] = useState(100);
    const [imageData, setImageData] = useState("");
    const [city, setCity] = useState("");
    const [address, setAddress] = useState("");
    const [latitude, setLatitude] = useState(0);
    const [longitude, setLongitude] = useState(0);
    const [startTime, setStartTime] = useState("");
    const [endTime, setEndTime] = useState("");

    const createEvent = async (e) => {
        e.preventDefault();
        try {
            await axios.get("/api/v1/auth/refresh");
            const pack = {
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
            console.log(pack);
            const res = await axios.post("/api/v1/events", pack);
            alert("Мероприятие создано!");
        } catch (error) {
            alert("Ошибка создания мероприятия!\n" + explainRequestError(error));
            console.error(error)
        }
    };

    const handleImageData = (e) => {
        const file = e.target.files[0];
        if (!file) return;

        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = () => {
            setImageData(reader.result);
        };
        reader.onerror = (error) => {
            alert("Невозможно загрузить изображение!\n"+error.message)
            console.error(error);
        };
    };

    return (
        <div className="create-event">
            <h1>Создать мероприятие</h1>
            <form onSubmit={createEvent} className="form">
                    <label for="name">Название:</label>
                    <input id="name" type="text" value={name} onChange={e=>setName(e.target.value)} required />
                    <label for="desctiption">Описание:</label>
                    <textarea id="description" type="text" value={description} onChange={e=>setDescription(e.target.value)} required />
                    <label for="category">Категория:</label>
                    <input id="category" type="text" value={category} onChange={e=>setCategory(e.target.value)} required />
                    <label for="maxParticipants">Участники:</label>
                    <input id="maxParticipants" type="number" min="1" value={maxParticipants} onChange={e=>setMaxParticipants(e.target.value)} required />
                    <label for="imageData">Иконка:</label>
                    <input id="imageData" type="file" accept="image/*" onChange={e=>handleImageData(e)} />
                    <label for="city">Город:</label>
                    <input id="city" type="text" value={city} onChange={e=>setCity(e.target.value)} />
                    <label for="address">Адрес:</label>
                    <input id="address" type="text" value={address} onChange={e=>setAddress(e.target.value)} />
                    <label for="latitude">Широта:</label>
                    <input id="latitude" type="number" min="-90" max="90" value={latitude} onChange={e=>setLatitude(e.target.value)} />
                    <label for="longitude">Долгота:</label>
                    <input id="longitude" type="number" value={longitude} min="-180" max="180" onChange={e=>setLongitude(e.target.value)} />
                    <label for="startTime">Время начала:</label>
                    <input id="startTime" type="datetime-local" value={startTime} onChange={e=>setStartTime(e.target.value)} required />
                    <label for="endTime">Время окончания:</label>
                    <input id="endTime" type="datetime-local" value={endTime} onChange={e=>setEndTime(e.target.value)} />
                <button type="submit">Создать мероприятие</button>
            </form>
        </div>
    )
}
