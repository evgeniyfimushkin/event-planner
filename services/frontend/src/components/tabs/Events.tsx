import "./Events.css";

import React from "react";
import { useContext, useEffect, useState } from "react";
import Grid from "../cards/Grid";
import FloatingButton from "../misc/FloatingButton";
import { useNavigate } from "react-router-dom";
import { createPortal } from "react-dom";
import ModalWindow from "../misc/ModalWindow";
import CreateEvent from "../event/CreateEvent";
import { authCall, explainRequestError } from "../../utilities/Utilities";
import AuthContext from "../../services/AuthContext";
import { Event, Registration } from "../../utilities/Types";
import { API } from "../../utilities/API";

export default function Events() {
    const [events, setEvents] = useState<Array<Event>>([]);
    const [subscriptions, setSubscriptions] = useState<Array<Registration>>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);
    const [showCreateEvent, setShowCreateEvent] = useState<boolean>(false);
    const { signOut } = useContext(AuthContext);
    const navigate = useNavigate();
    const [query, setQuery] = useState<string>("");
    const [includePrevious, setIncludePrevious] = useState<boolean>(false);

    const fetchData = async (e?) => {
        setLoading(true);
        setError(null);
        try {
            await authCall(async () => { // success
                setLoading(true);
                let events;
                if (query) {
                    const responseEvents = await API.Events.Search(query);
                    setEvents(responseEvents.data);
                } else {
                    const responseEvents = await API.Events.Get();
                    const responseEventsPrevious = includePrevious ? await API.Events.GetPrevious() : {data: []};
                    setEvents([...responseEventsPrevious.data, ...responseEvents.data]);
                }
                const responseSubscriptions = await API.Registrations.GetMy();
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
    const search = async (e) => {
        e.preventDefault();
        fetchData();
    }

    useEffect(() => {
        fetchData();
    }, [includePrevious]);

    return (
        <>
            <form className="search" onSubmit={search}>
                <input type="text" placeholder="Запрос" value={query} onChange={(e) => setQuery(e.target.value)} />
                <button type="submit">Поиск</button>
            </form>
            {!query && <>
                <div className="options">
                    <label>
                        Включить предыдущие мероприятия
                        <input id="includePrevious" type="checkbox" checked={includePrevious} onChange={e => {setIncludePrevious(e.target.checked)}} />
                    </label>
                </div>
            </>}
            {
                (loading && <p>Загрузка...</p>) ||
                (error && <p>{error}</p> ) ||
                <>
                    <Grid cards={events} subscriptions={subscriptions} />
                    <FloatingButton text="Добавить мероприятие" onClick={()=>setShowCreateEvent(true)} />
                    {showCreateEvent && createPortal(
                        <ModalWindow onClose={()=>{setShowCreateEvent(false);fetchData()}}>
                            <CreateEvent />
                        </ModalWindow>, document.body
                    )}
                </>
            }
        </>
    );
}
