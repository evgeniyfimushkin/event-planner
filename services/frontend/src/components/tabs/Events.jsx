import { useContext, useEffect, useState } from "react";
import axios from "axios";
import Grid from "../cards/Grid";
import FloatingButton from "../misc/FloatingButton";
import { useNavigate } from "react-router-dom";
import { createPortal } from "react-dom";
import ModalWindow from "../misc/ModalWindow";
import CreateEvent from "../event/CreateEvent";
import { authCall } from "../../services/Utilities";
import AuthContext from "../../services/AuthContext";

export default function Events() {
    const [events, setEvents] = useState([]);
    const [subscriptions, setSubscriptions] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [showCreateEvent, setShowCreateEvent] = useState(false);
    const { logout } = useContext(AuthContext);
    const navigate = useNavigate();

    const fetchData = async () => {
        try {
            await authCall(async () => { // success
                setLoading(true);
                const responseEvents = await axios.get("/api/v1/events");
                const responseSubscriptions = await axios.get("/api/v1/registrations/my");
                setEvents(responseEvents.data);
                setSubscriptions(responseSubscriptions.data);
                // console.log(response.data);
            }, (err) => { // unauthorized
                logout();
                navigate("/login");
            });
        } catch (err) {
            setError("Не удалось загрузить мероприятия");
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
            <Grid cards={events} subscriptions={subscriptions} />
            <FloatingButton text="Добавить мероприятие" onClick={()=>setShowCreateEvent(true)} />
            {showCreateEvent && createPortal(
                <ModalWindow onClose={()=>{setShowCreateEvent(false);fetchData()}}>
                    <CreateEvent />
                </ModalWindow>, document.body
            )}
        </>
    );
}
