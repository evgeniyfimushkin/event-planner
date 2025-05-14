import axios from "axios"
import { Event, Review } from "./Types";

interface RequestArguments {
    [key: string]: any
}

interface LoginCredentials {
    username: string,
    passhash: string,
}

interface RegisterCredentials {
    username: string,
    email: string,
    passhash: string,
}

export const API = {
    Auth: {
        Refresh: async () => axios.get("/api/v1/auth/refresh"),
        Login: async (args: LoginCredentials) => axios.post("/api/v1/auth/login", args),
        Register: async (args: RegisterCredentials) => axios.post("/api/v1/auth/register", args),
    },
    Events: {
        Create: async (args: Event) => axios.post("/api/v1/events", args),
        Get: async () => axios.get("/api/v1/events"),
        GetPrevious: async () => axios.get("/api/v1/events/previous"),
        Search: async (query: string) => axios.get(`/api/v1/events/search?${query}`),
        Update: async (args: Event) => {console.log("API.Events.Update", args)} // async (args: Event) => axios.post("/api/v1/events", args), // todo use actual endpoint
    },
    Registrations: {
        Create: async (id: number) => axios.post("/api/v1/registrations", { event_id: id }),
        Delete: async (id: number) => axios.delete("/api/v1/registrations", { data: {event_id: id} }),
        GetMy: async () => axios.get("/api/v1/registrations/my"),
    },
    Reviews: {
        Get: async (event_id: number) => {
            const g = (
                id: number,
                user_id: number,
                content: string,
                posted_at: string,
                updated_at: string,
            ) => ({id, user_id, content, posted_at, updated_at})
            return {data: [
                g(1, 1, "первый!11!!", "2025-05-21T21:04:00Z", "2025-05-27T23:06:00Z"),
                g(2, 2, "второй", "2025-05-21T21:04:00Z", "2025-05-27T23:06:00Z"),
                g(3, 3, `г${"о".repeat(200)}л`, "2025-05-21T21:04:00Z", "2025-05-27T23:06:00Z"),
                g(4, 4, `ты совсем? мы фронтенд тестируем`, "2025-05-21T21:04:00Z", "2025-05-27T23:06:00Z"),
            ] as any}
        }, // todo use actual endpoint
    }
};
