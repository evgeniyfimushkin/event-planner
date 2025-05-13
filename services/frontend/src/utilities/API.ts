import axios from "axios"
import { Event } from "./Types";

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
};
