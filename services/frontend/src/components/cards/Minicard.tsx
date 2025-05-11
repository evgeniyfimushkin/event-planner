import "./Cards.css"

import React, { useCallback } from "react";
import { useState } from "react";
import { localDate } from "../../utilities/Utilities";
import { Event } from "../../utilities/Types";

export default function Minicard({event, callback}: {
    event: Event,
    callback?: (event: Event) => any
}) {
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

    const onClickCallback = callback ? (() => callback(event)) : (()=>{});

    return (
        <div className="card mini" onClick={onClickCallback}>
            <div className="line">
                {image_data && <img src={image_data}/>}
                <h1 className="title">{name}</h1>
                {description && <p className="description">{description}</p>}
            </div>
            {fullAddress && (
                <p>Местоположение: {fullAddress}</p>
            )}
            {start_time && <p className="startTime">Начало: {localDate(new Date(start_time))}</p>}
        </div>
    )
}
