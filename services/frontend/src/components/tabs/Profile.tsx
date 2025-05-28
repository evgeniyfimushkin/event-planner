import "./Profile.css";

import React from "react";
import { useContext, useEffect, useState } from "react";
import { data, useNavigate } from "react-router-dom";
import { authCall, explainRequestError } from "../../utilities/Utilities";
import AuthContext from "../../services/AuthContext";
import { UserData, UserSettings } from "../../utilities/Types";
import { API } from "../../utilities/API";
import Menu from "../menu/Menu";
import CryptoJS from "crypto-js";

export default function Profile() {
    const [settings, setSettings] = useState<UserSettings | null>(null);
    const [loadingSettings, setLoadingSettings] = useState<boolean>(true);
    const [errorSettings, setErrorSettings] = useState<string | null>(null);
    const [image, setImage] = useState<string | null>(null);
    const { signOut, username } = useContext(AuthContext);
    const navigate = useNavigate();

    const fetchData = async (e?) => {
        setLoadingSettings(true);
        setErrorSettings(null);
        try {
            await authCall(async () => { // success
                setLoadingSettings(true);
                const responseSettings = await API.Users.GetSettingsMy();
                setSettings(responseSettings.data);
                setImage(responseSettings.data.picture ?? null);
            }, (err) => { // unauthorized
                signOut();
                navigate("/signIn");
            });
        } catch (err) {
            setErrorSettings("Не удалось загрузить данные о пользователе!\n"+explainRequestError(err));
            console.error(err);
        } finally {
            setLoadingSettings(false);
        }
    }
    useEffect(() => {
        fetchData();
    }, []);

    const updateData = async (e) => {
        e.preventDefault();
        try {
            await API.Auth.Refresh();
            await API.Users.UpdateSettingsMy(settings!);
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
            const newSettings = settings!;
            newSettings.picture = reader.result as string;
            setSettings(newSettings);
            setImage(newSettings.picture);
            console.log("set image")
        };
        reader.onerror = (pe) => {
            alert("Невозможно загрузить изображение!\n"+reader.error)
            console.error(reader.error);
        };
    };
    const setTelegram = (e) => {
        const newSettings = settings!;
        if (e !== "") newSettings.telegram = e;
        setSettings(newSettings);
    }
    const setPasshash = (e) => {
        const newSettings = settings!;
        if (e !== "") newSettings.passhash = CryptoJS.SHA256(e).toString(CryptoJS.enc.Hex);
        setSettings(newSettings);
    }

    return (
        <>{
            (loadingSettings && <p>Загрузка...</p>) ||
            (errorSettings && <p>{errorSettings}</p> ) ||
            <>
            <div className="profile">
                <h1>Добро пожаловать,<br/>{username}</h1>
                {image && <img src={image} onClick={imageClickHandler}/>}
                <form onSubmit={updateData} className="form">
                    {image ? <></> : <label htmlFor="imageData">Иконка:</label>}
                    <input id="imageData" type="file" accept="image/*" onChange={e=>handleImageData(e)} className={image ? "hidden" : ""} />
                    {/* <p><label>Интервал уведомления (с): <input id="interval" type="number" min="1" value={settings!.interval} onChange={e=>setInterval(+e.target.value)} required /></label></p> */}
                    <p>Сменить Telegram: <input type="text" placeholder="Telegram" value={settings!.telegram} onChange={(e) => setTelegram(e.target.value)} /></p>
                    <p>Сменить пароль: <input type="password" placeholder="Пароль" value={settings!.passhash} onChange={(e) => setPasshash(e.target.value)} /></p>
                    <p><button type="submit">Применить изменения</button></p>
                </form>
            </div>
            </>
        }</>
    );
}
