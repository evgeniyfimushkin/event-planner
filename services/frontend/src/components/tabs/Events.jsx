import { useEffect, useState } from "react";
import axios from "axios";
import Grid from "../cards/Grid";
import FloatingButton from "../misc/FloatingButton";
import { useNavigate } from "react-router-dom";

export default function Events() {
    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    const navigate = useNavigate();

    useEffect(() => {
        const fetchData = async () => {
            try {
                await axios.get("/api/v1/auth/refresh");
                const response = await axios.get("/api/v1/events");
                setData(response.data);
                console.log(response.data);
            } catch (err) {
                setError("Can't receive events");
                console.error(err);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, []);

    if (loading) return <p>Loading...</p>;
    if (error) return <p>{error}</p>;

    return (
        <>
            <Grid cards={data} />
            <FloatingButton text="Добавить мероприятие" target={"/events"} />
        </>
    );
}
