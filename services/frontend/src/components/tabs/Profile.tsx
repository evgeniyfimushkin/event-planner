import "./Profile.css";

import React from "react";
import { useContext, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { authCall, explainRequestError } from "../../utilities/Utilities";
import AuthContext from "../../services/AuthContext";
import { UserData } from "../../utilities/Types";
import { API } from "../../utilities/API";
import Menu from "../menu/Menu";

export default function Profile() {
    const [data, setData] = useState<UserData | null>(null);
    const [image, setImage] = useState<string | null>(null);
    const [loadingData, setLoadingData] = useState<boolean>(true);
    const [errorData, setErrorData] = useState<string | null>(null);
    const { signOut } = useContext(AuthContext);
    const navigate = useNavigate();

    const fetchData = async (e?) => {
        setLoadingData(true);
        setErrorData(null);
        try {
            await authCall(async () => { // success
                setLoadingData(true);
                const responseData = await API.Users.GetMy();
                setData(responseData.data);
                setImage(responseData.data.picture ?? null);
            }, (err) => { // unauthorized
                signOut();
                navigate("/signIn");
            });
        } catch (err) {
            setErrorData("Не удалось загрузить информацию о пользователе!\n"+explainRequestError(err));
            console.error(err);
        } finally {
            setLoadingData(false);
        }
    }
    useEffect(() => {
        fetchData();
    }, []);

    const updateData = async (e) => {
        e.preventDefault();
        try {
            await API.Auth.Refresh();
            await API.Users.UpdateMy(data!);
            alert("Профиль обновлён!");
        } catch (error) {
            alert("Ошибка обновления!\n" + explainRequestError(error));
            console.error(error)
        }
    };
    
    const imageClickHandler = (e) => {
        e.preventDefault();
        const imageInput = document.getElementById("imageData");
        imageInput!.click();
    }
    const handleImageData = (e) => {
        console.log("handling image...")
        const file = e.target.files[0];
        if (!file) return;
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = () => {
            const newData = data!;
            newData.picture = reader.result as string;
            setData(newData);
            setImage(newData.picture);
            console.log("set image")
        };
        reader.onerror = (pe) => {
            alert("Невозможно загрузить изображение!\n"+reader.error)
            console.error(reader.error);
        };
    };

    return (
        <>{
            (loadingData && <p>Загрузка...</p>) ||
            (errorData && <p>{errorData}</p> ) ||
            <>
            <div className="profile">
                <h1>Добро пожаловать,<br/>{data!.username}</h1>
                {image && <img src={image} onClick={imageClickHandler}/>}
                <form onSubmit={updateData} className="form">
                    {image ? <></> : <label htmlFor="imageData">Иконка:</label>}
                    <input id="imageData" type="file" accept="image/*" onChange={e=>handleImageData(e)} className={image ? "hidden" : ""} />
                    <p><button type="submit">Применить изменения</button></p>
                </form>
            </div>
            </>
        }</>
    );
}
