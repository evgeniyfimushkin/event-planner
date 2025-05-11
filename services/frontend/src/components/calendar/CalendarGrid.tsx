import "./Calendar.css";

import React, { useState } from "react";
import Minicard from "../cards/Minicard";
import { Event } from "../../utilities/Types";
import { createPortal } from "react-dom";
import ModalWindow from "../misc/ModalWindow";
import Card from "../cards/Card";

export default function CalendarGrid({events=[]}: {
    events: Array<Event>
}) {
    const getCurrentDate = () => {
        return new Date();
    }
    const isEventNow = (event: Event, day: number, hour: number) => {
        const eventDate = new Date(event.start_time);
        const date = getCurrentDate();
        const start = new Date(date.getFullYear(), date.getMonth(), day, hour);
        const end = new Date(date.getFullYear(), date.getMonth(), day, hour+1);
        const isNow = eventDate > start && eventDate < end;
        return isNow
    }
    const isToday = (day: number) => {
        const date = getCurrentDate();
        return date.getDate() === day;
    }
    const isNow = (hour: number) => {
        const date = getCurrentDate();
        return date.getHours() === hour;
    }

    const currentDate = getCurrentDate();
    const firstDayDate = new Date(currentDate.getFullYear(), currentDate.getMonth(), 1);
    const lastDayDate = new Date(currentDate.getFullYear(), currentDate.getMonth()+1, 0);
    const firstHour = 0;
    const lastHour = 23;
    const rangeDate: Array<number> = [];
    for (let day = firstDayDate.getDate(); day <= lastDayDate.getDate(); day++) {
        rangeDate.push(day);
    }
    const rangeHour: Array<number> = [];
    for (let hour = firstHour; hour <= lastHour; hour++) {
        rangeHour.push(hour);
    }

    const [showModal, setShowModal] = useState<boolean>(false);
    const [event, setEvent] = useState<Event | null>(null);

    const clickCallback = async (event: Event) => {
        setEvent(event);
        setShowModal(true);
    }

    return (
        <>
            <div className="calendar">
                <table className="calendar">
                    <caption>{getCurrentDate().toLocaleString("ru", {month: "long", year: "numeric"})}</caption>
                    <thead>
                        <tr>
                            <td>\</td>
                            {rangeDate.map((day,di)=><td key={di}>{day}</td>)}
                        </tr>
                    </thead>
                    <tbody>
                        {rangeHour.map((hour,hi) => <tr key={hi}>
                            <td>{hour}:00</td>
                            {rangeDate.map((day,di) => <td key={di} className={(isToday(day) && "today" || "") + " " + (isNow(hour) && "now" || "")}>
                                {events.filter(e=>isEventNow(e,day,hour)).map(e=><Minicard event={e} callback={clickCallback} />)}
                            </td>)}
                        </tr>)}
                    </tbody>
                </table>
            </div>
            {showModal && createPortal(
                <ModalWindow onClose={()=>{setShowModal(false);}}>
                    <Card event={event!} />
                </ModalWindow>,
                document.body
            )}
        </>
    );
}
