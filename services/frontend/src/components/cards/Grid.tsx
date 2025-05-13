import "./Cards.css"

import React, { useState } from "react"
import Minicard from "./Minicard.jsx"
import { Event, Registration } from "../../utilities/Types"
import { createPortal } from "react-dom"
import ModalWindow from "../misc/ModalWindow.js"
import Card from "./Card.js"
import EditEvent from "../event/EditEvent.js"

export default function Grid({cards=[], subscriptions=[]}: {
    cards: Array<Event>,
    subscriptions: Array<Registration>
}) {
    const [showCard, setShowCard] = useState<boolean>(false);
    const [showEdit, setShowEdit] = useState<boolean>(false);
    const [event, setEvent] = useState<Event | null>(null);

    const cardClickCallback = async (event: Event) => {
        setEvent(event);
        setShowCard(true);
    }

    const editOpenCallback = async () => {
        setShowEdit(true);
    }

    return (
        <>
            <div className="container">
                <div className="grid">
                    {cards.map((card,index) => (
                        <Minicard key={card.id} event={card} callback={cardClickCallback} />
                    ))}
                </div>
            </div>
            {showCard && createPortal(
                <ModalWindow buttons={[
                    {name: "Редактировать", onClick: editOpenCallback},
                    {name: "Закрыть", onClick: ()=>{setShowCard(false);}},
                ]}>
                    <Card event={event!} />
                </ModalWindow>,
                document.body
            )}
            {showEdit && createPortal(
                <ModalWindow buttons={[
                    {name: "Закрыть", onClick: ()=>{setShowEdit(false);}}, // todo refresh
                ]}>
                    <EditEvent event={event!} />
                </ModalWindow>,
                document.body
            )}
        </>
    )
}
