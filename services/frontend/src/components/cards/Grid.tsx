import "./Cards.css"

import React, { useState } from "react"
import Minicard from "./Minicard.jsx"
import { Event, Registration } from "../../utilities/Types"
import { createPortal } from "react-dom"
import ModalWindow from "../misc/ModalWindow.js"
import Card from "./Card.js"

export default function Grid({cards=[], subscriptions=[]}: {
    cards: Array<Event>,
    subscriptions: Array<Registration>
}) {
    const [showModal, setShowModal] = useState<boolean>(false);
    const [event, setEvent] = useState<Event | null>(null);

    const clickCallback = async (event: Event) => {
        setEvent(event);
        setShowModal(true);
    }

    return (
        <>
            <div className="container">
                <div className="grid">
                    {cards.map((card,index) => (
                        <Minicard key={card.id} event={card} callback={clickCallback} />
                    ))}
                </div>
            </div>
            {showModal && createPortal(
                <ModalWindow onClose={()=>{setShowModal(false);}}>
                    <Card event={event!} />
                </ModalWindow>,
                document.body
            )}
        </>
    )
}
