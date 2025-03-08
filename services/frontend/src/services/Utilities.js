import axios from "axios";

export function localDate(date) {
    return date.toLocaleString("ru", {
        dateStyle: "medium",
        timeStyle: "long",
    });
}

export async function authCall(request, handler401) {
    const handlerOther = (err) => {
        throw new Error("Auth call error (not 401)", {cause: err});
    }
    try {
        const out = await request();
        return out;
    } catch (err) {
        if (err.status === 401) {
            try {
                await axios.get("/api/v1/auth/refresh");
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

export function explainRequestError(error) {
    return (error.response)
        ? `Статус ответа: ${error.response.status}\nСообщение:\n${error.response.data}`
        : `Сообщение библиотеки:\n${error.message}`;
}
