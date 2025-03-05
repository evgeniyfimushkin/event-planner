import "./Event.css"

import { useState, useContext } from "react";
import axios from "axios";

export default function CreateEvent({}) {
    const [name, setName] = useState("Untitled");
    const [description, setDescription] = useState("No description");
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
            const pack = {
                name,
                description,
                category,
                max_participants: maxParticipants,
                image_data: imageData,
                city,
                address,
                latitude,
                longitude,
                // start_time: startTime,
                // end_time: endTime,
            };
            console.log(pack);
            const res = await axios.post("http://localhost/api/v1/events", pack);
            alert("Event created!");
        } catch (error) {
            console.error("Error adding event",e);
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
            console.error("Unable to load image!", error);
        };
    };

    return (
        <div>
            <h1>Create event</h1>
            <form onSubmit={createEvent} className="form">
                <div>
                    <label for="name">Name:</label>
                    <input id="name" type="text" value={name} onChange={e=>setName(e.target.value)} />
                </div>
                <div>
                    <label for="desctiption">Description:</label>
                    <input id="description" type="text" value={description} onChange={e=>setDescription(e.target.value)} />
                </div>
                <div>
                    <label for="category">Category:</label>
                    <input id="category" type="text" value={category} onChange={e=>setCategory(e.target.value)} />
                </div>
                <div>
                    <label for="maxParticipants">Participants:</label>
                    <input id="maxParticipants" type="number" min="1" value={maxParticipants} onChange={e=>setMaxParticipants(e.target.value)} />
                </div>
                <div>
                    <label for="imageData">Icon:</label>
                    <input id="imageData" type="file" accept="image/*" onChange={e=>handleImageData(e)} />
                </div>
                <div>
                    <label for="city">City:</label>
                    <input id="city" type="text" value={city} onChange={e=>setCity(e.target.value)} />
                </div>
                <div>
                    <label for="address">Address:</label>
                    <input id="address" type="text" value={address} onChange={e=>setAddress(e.target.value)} />
                </div>
                <div>
                    <label for="latitude">Latitude:</label>
                    <input id="latitude" type="number" min="-90" max="90" value={latitude} onChange={e=>setLatitude(e.target.value)} />
                </div>
                <div>
                    <label for="longitude">Longitude:</label>
                    <input id="longitude" type="number" value={longitude} min="-180" max="180" onChange={e=>setLongitude(e.target.value)} />
                </div>
                {/* <div>
                    <label for="startTime">Start time:</label>
                    <input id="startTime" type="datetime-local" value={startTime} onChange={e=>setStartTime(e.target.value)} />
                </div>
                <div>
                    <label for="endTime">End time:</label>
                    <input id="endTime" type="datetime-local" value={endTime} onChange={e=>setEndTime(e.target.value)} />
                </div> */}
                <button type="submit">Add event</button>
            </form>
        </div>
    )
}
