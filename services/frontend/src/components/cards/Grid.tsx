import "./Cards.css"

import React from "react"
import Minicard from "./Minicard.jsx"
import { Event, Registration } from "../../utilities/Types"

export default function Grid({cards=[], subscriptions=[]}: {
    cards: Array<Event>,
    subscriptions: Array<Registration>
}) {
    return (
        <div className="container">
            <div className="grid">
                {cards.map((card,index) => (
                    <Minicard key={card.id} event={card} />
                ))}
            </div>
        </div>
    )
}
