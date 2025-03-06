import Minicard from "./Minicard.jsx"
import "./Cards.css"

export default function Grid({cards=[], subscriptions=[]}) {
    return (
        <div className="container">
            <div className="grid">
                {cards.map((card,index) => (
                    <Minicard key={card.id} event={card} subscribed={subscriptions.some(c=>c.event_id === card.id)}/>
                ))}
            </div>
        </div>
    )
}
