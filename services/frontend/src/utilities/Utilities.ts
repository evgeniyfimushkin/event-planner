import axios, { AxiosError } from "axios";
import { API } from "./API";

export function localeDateString(date: Date): string { // для отображения
    return date.toLocaleString("ru", {
        dateStyle: "medium",
        timeStyle: "long",
    });
}

export function toLocalDate(date: Date): string { // для datetime-local
    const offsetDate = new Date(date.getTime() - date.getTimezoneOffset() * 60000);
    const formattedDate = offsetDate.toISOString();
    const formattedLocalDate = formattedDate.replace(/:\d{2}\.\d{3}Z$/, '');
    return formattedLocalDate;
}

export async function authCall(request: (...args: any[]) => any, handler401: (...args: any[]) => any) {
    const handlerOther = (err) => {
        throw new Error("Auth call error (not 401)", {cause: err});
    }
    try {
        const out = await request();
        return out;
    } catch (err) {
        if (err.status === 401) {
            try {
                await API.Auth.Refresh();
                const out = await request();
                return out;
            } catch (err) {
                if (err.status === 401) {
                    handler401(err);
                    return;
                }
                handlerOther(err);
                return;
            }
        }
        handlerOther(err);
        return;
    }
};

export function explainRequestError(error: Error | AxiosError) {
    return (axios.isAxiosError(error))
        ? `Статус ответа: ${error.response?.status}\nСообщение:\n${error.response?.data}`
        : `Сообщение библиотеки:\n${error.message}`;
}
