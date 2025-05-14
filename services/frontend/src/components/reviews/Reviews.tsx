import "./Reviews.css"

import React, { useContext } from "react";
import { useState } from "react";
import { explainRequestError, localeDateString } from "../../utilities/Utilities";
import { useEffect } from "react";
import { authCall } from "../../utilities/Utilities";
import { Review } from "../../utilities/Types";
import { useNavigate } from "react-router-dom";
import AuthContext from "../../services/AuthContext";
import { API } from "../../utilities/API";
import ReviewElement from "./Review";

export default function Reviews({event_id}: {
    event_id: number
}) {
    const [reviews, setReviews] = useState<Array<Review>>([]);
    const [loadingReviews, setLoadingReviews] = useState<boolean>(true);
    const [errorReviews, setErrorReviews] = useState<string | null>(null);
    const { signOut } = useContext(AuthContext);
    const navigate = useNavigate();

    const fetchReviews = async () => {
        try {
            await authCall(async () => { // success
                setLoadingReviews(true);
                const responseReviews = await API.Reviews.Get(event_id);
                setReviews(responseReviews.data);
            }, (err) => { // unauthorized
                signOut();
                navigate("/signIn");
            });
        } catch (err) {
            setErrorReviews("Не удалось загрузить отзывы!\n"+explainRequestError(err));
            console.error(err);
        } finally {
            setLoadingReviews(false);
        }
    }
    useEffect(() => {
        fetchReviews();
    }, []);

    return (
        <>{
            (loadingReviews && <p>Загрузка...</p>) ||
            (errorReviews && <p>{errorReviews}</p> ) ||
            <>
            <div className="reviews-list">
                {reviews.map((review,index)=><ReviewElement key={review.id} review={review} />)}
            </div>
            </>
        }</>
    )
}
