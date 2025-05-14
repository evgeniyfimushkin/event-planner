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
import Card from "../cards/Card";
import EditEvent from "../event/EditEvent";

export default function Events() {
    const [events, setEvents] = useState<Array<Event>>([]);
    const [subscriptions, setSubscriptions] = useState<Array<Registration>>([]);
    const [loadingEvents, setLoadingEvents] = useState<boolean>(true);
    const [errorEvents, setErrorEvents] = useState<string | null>(null);
    const [showCreateEvent, setShowCreateEvent] = useState<boolean>(false);
    const { signOut } = useContext(AuthContext);
    const navigate = useNavigate();
    const [query, setQuery] = useState<string>("");
    const [includePrevious, setIncludePrevious] = useState<boolean>(false);

    const [showCard, setShowCard] = useState<boolean>(false);
    const [showEditEvent, setShowEditEvent] = useState<boolean>(false);
    const [selectedEvent, setSelectedEvent] = useState<Event | null>(null);

    const fetchEvents = async (e?) => {
        setLoadingEvents(true);
        setErrorEvents(null);
        try {
            await authCall(async () => { // success
                setLoadingEvents(true);
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
            setErrorEvents("Не удалось загрузить мероприятия!\n"+explainRequestError(err));
            console.error(err);
        } finally {
            setLoadingEvents(false);
        }
    }
    const search = async (e) => {
        e.preventDefault();
        fetchEvents();
    }
    useEffect(() => {
        fetchEvents();
    }, [includePrevious]);

    const cardClickCallback = async (event: Event) => {
        setSelectedEvent(event);
        setShowCard(true);
    }
    const editOpenCallback = async () => {
        setShowEditEvent(true);
    }

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
                (loadingEvents && <p>Загрузка...</p>) ||
                (errorEvents && <p>{errorEvents}</p> ) ||
                <>
                    <Grid cards={events} subscriptions={subscriptions} cardCallback={cardClickCallback} />
                    <FloatingButton text="Добавить мероприятие" onClick={()=>setShowCreateEvent(true)} />
                    {showCreateEvent && createPortal(
                        <ModalWindow buttons={[
                            {name: "Закрыть", onClick: ()=>{setShowCreateEvent(false);fetchEvents()}}
                        ]}>
                            <CreateEvent />
                        </ModalWindow>, document.body
                    )}
                </>
            }
            {showCard && createPortal(
                <ModalWindow buttons={[
                    {name: "Редактировать", onClick: editOpenCallback},
                    {name: "Закрыть", onClick: ()=>{setShowCard(false);}},
                ]}>
                    <Card event={selectedEvent!} />
                </ModalWindow>,
                document.body
            )}
            {showEditEvent && createPortal(
                <ModalWindow buttons={[
                    {name: "Закрыть", onClick: ()=>{setShowEditEvent(false);}}, // todo refresh
                ]}>
                    <EditEvent event={selectedEvent!} />
                </ModalWindow>,
                document.body
            )}
        </>
    );
}
