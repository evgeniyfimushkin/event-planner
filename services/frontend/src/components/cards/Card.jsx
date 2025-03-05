import "./Cards.css"

export default function Card({event}) {
    const {
        name,
        description,
        category,
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
    return (
        <div className="card">
            <div className="line">
                {image_data && <img src={image_data || null}/>}
                <h1 className="title">{name}</h1>
            </div>
            {description && <p className="description">{description}</p>}
            {max_participants && <p className="maxParticipants">{max_participants} мест</p>}
            {fullAddress && (
                <p>Местоположение: {fullAddress}</p>
            )}
            {start_time && <p className="startTime">Начало: {start_time}</p>}
            {end_time && <p className="endTime">Окончание: {end_time}</p>}
            {category && <p className="category">{category}</p>}
        </div>
    )
}
