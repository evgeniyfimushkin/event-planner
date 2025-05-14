import "./Cards.css"

import React from "react"
import Minicard from "./Minicard.jsx"
import { Event, Registration } from "../../utilities/Types"

export default function Grid({cards=[], subscriptions=[], cardCallback}: {
    cards: Array<Event>,
    subscriptions: Array<Registration>,
    cardCallback?: (event: Event) => any,
}) {
    const cardClickCallback = cardCallback ? ((event) => cardCallback(event)) : (()=>{});

    return (
        <>
            <div className="container">
                <div className="grid">
                    {cards.map((card,index) => (
                        <Minicard key={card.id} event={card} callback={cardClickCallback} />
                    ))}
                </div>
            </div>
        </>
    )
}
