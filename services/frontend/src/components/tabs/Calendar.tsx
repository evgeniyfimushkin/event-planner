import React from "react";

import { useContext, useEffect, useState } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { createPortal } from "react-dom";
import ModalWindow from "../misc/ModalWindow.js";
import { authCall, explainRequestError } from "../../utilities/Utilities.js";
import AuthContext from "../../services/AuthContext.jsx";
import CalendarGrid from "../calendar/CalendarGrid.tsx";
import { Event, Registration } from "../../utilities/Types.ts";

export default function Calendar() {
    const [events, setEvents] = useState<Array<Event>>([]);
    const [subscriptions, setSubscriptions] = useState<Array<Registration>>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<null | string>(null);
    const { signOut } = useContext(AuthContext);
    const navigate = useNavigate();

    const fetchData = async () => {
        try {
            await authCall(async () => { // success
                setLoading(true);
                const responseEvents = await axios.get("/api/v1/events");
                const responseSubscriptions = await axios.get("/api/v1/registrations/my");
                setEvents(responseEvents.data);
                setSubscriptions(responseSubscriptions.data);
            }, (err) => { // unauthorized
                signOut();
                navigate("/signIn");
            });
        } catch (err) {
            setError("Не удалось загрузить мероприятия!\n"+explainRequestError(err));
            console.error(err);
        } finally {
            setLoading(false);
        }
    }

    useEffect(() => {
        fetchData();
    }, []);

    if (loading) return <p>Загрузка...</p>;
    if (error) return <p>{error}</p>;

    return (
        <>
            <CalendarGrid events={events.filter(e => subscriptions.some(s => s.event_id === e.id))} />
            {/* <Grid cards={events} subscriptions={subscriptions} />
            <FloatingButton text="Добавить мероприятие" onClick={()=>setShowCreateEvent(true)} /> */}
            {/* {showCreateEvent && createPortal(
                <ModalWindow onClose={()=>{setShowCreateEvent(false);fetchData()}}>
                    <CreateEvent />
                </ModalWindow>, document.body
            )} */}
        </>
    );
}
