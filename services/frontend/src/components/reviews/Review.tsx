import "./Reviews.css"

import React from "react";
import { Review as ReviewType } from "../../utilities/Types";

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

    return (
        <div className="review">
            <p className="header">Пользователь <span className="username">{user_id}</span></p>
            <p className="review-text">{content}</p>
        </div>
    )
}
