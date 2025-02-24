import Card from "./Card.jsx"
import "./Cards.css"

export default function Grid({cards=[{title: "Мероприятие 1"},{description: "Описание мероприятия 2"}, {}, {}]}) {
    return (
        <div className="container">
            <div className="grid">
                {cards.map((card,index) => (
                    <Card key={index} title={card.title} description={card.description} time={card.time} />
                ))}
            </div>
        </div>
    )
}
