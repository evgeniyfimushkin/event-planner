import "./Reviews.css"

import React, { useEffect, useState } from "react";
import { Review as ReviewType, UserData } from "../../utilities/Types";
import { authCall, explainRequestError } from "../../utilities/Utilities";
import { API } from "../../utilities/API";

export default function Review({review}: {
    review: ReviewType
}) {
    const {
        id,
        user_id,
        content,
        posted_at,
        updated_at,
    } = review;

    const [data, setData] = useState<UserData | null>(null);
    const [loadingData, setLoadingData] = useState<boolean>(true);
    const [errorData, setErrorData] = useState<string | null>(null);

    const fetchData = async () => {
        try {
            await authCall(async () => { // success
                setLoadingData(true);
                const responseData = await API.Users.Get(user_id);
                setData(responseData.data);
            }, (err) => { // unauthorized
                // signOut();
                // navigate("/signIn");
                throw new Error("Unauthorized"); // is this applicable?
            });
        } catch (err) {
            setErrorData("Не удалось загрузить данные пользователя!\n"+explainRequestError(err));
            console.error(err);
        } finally {
            setLoadingData(false);
        }
    }
    useEffect(() => {
        fetchData();
    }, []);

    return (
        <div className="review">
            <p className="header">{
                (loadingData && <>...</>) ||
                (errorData && <>Пользователь <span className="username">{user_id}</span></> ) ||
                <>
                    {data!.picture && <img src={data!.picture}/>}
                    <span className="username">{data!.username}</span>
                </>
            }</p>
            <p className="review-text">{content}</p>
        </div>
    )
}
