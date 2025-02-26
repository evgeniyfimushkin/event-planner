import "./Cards.css"

export default function Card({title="Без названия", description="Нет описания", time="Время не назначено"}) {
    return (
        <div className="card">
            <h1 className="title">{title}</h1>
            <p className="description">{description}</p>
            <p className="time">{time}</p>
        </div>
    )
}
