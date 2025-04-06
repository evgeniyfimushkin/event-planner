import axios, { AxiosError } from "axios";

export function localDate(date: Date): string {
    return date.toLocaleString("ru", {
        dateStyle: "medium",
        timeStyle: "long",
    });
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

export function explainRequestError(error: Error | AxiosError) {
    return (axios.isAxiosError(error))
        ? `Статус ответа: ${error.response?.status}\nСообщение:\n${error.response?.data}`
        : `Сообщение библиотеки:\n${error.message}`;
}
