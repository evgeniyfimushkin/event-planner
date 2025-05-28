import axios from "axios"
import { Event, Review, UserCredentials, UserData, UserSettings } from "./Types";
import { Picture } from "../assets/Placeholders";
import { useContext } from "react";
import AuthContext from "../services/AuthContext";

interface RequestArguments {
    [key: string]: any
}

interface TemporaryRegisterCredentials extends UserCredentials {
    email: string,
}

export const API = {
    Auth: {
        Refresh: async () => axios.get("/api/v1/auth/refresh"),
        Login: async (args: UserCredentials) => {
            const res = await axios.post("/api/v1/auth/login", args);
            return res;
        },
        Register: async (args: TemporaryRegisterCredentials) => {
            axios.post("/api/v1/auth/register", args)
        },
    },
    Events: {
        Create: async (args: Event) => axios.post("/api/v1/events", args),
        Get: async () => axios.get("/api/v1/events"),
        GetPrevious: async () => axios.get("/api/v1/events/previous"),
        Search: async (query: string) => axios.get(`/api/v1/events/search?${query}`),
        Update: async (args: Event) => {console.log("API.Events.Update", args)}, // async (args: Event) => axios.post("/api/v1/events", args), // todo use actual endpoint
        GetOwn: async () => {}, // todo
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
                g(3, 3, `третий отзыв`, "2025-05-21T21:04:00Z", "2025-05-27T23:06:00Z"),
                g(4, 4, `ты совсем? мы фронтенд тестируем`, "2025-05-21T21:04:00Z", "2025-05-27T23:06:00Z"),
            ] as any}
        }, // todo use actual endpoint
        Create: async (args: Review) => {
            console.log("API.Reviews.Create", args);
        }, // todo use actual endpoint
        Update: async (args: Review) => {
            console.log("API.Reviews.Update", args);
        }, // todo use actual endpoint
        Delete: async (id: number) => {
            console.log("API.Reviews.Delete", id);
        }
    },
    Users: {
        GetUsername: async () => {
            const { username } = useContext(AuthContext);
            return username;
        },
        GetDataMy: async () => {
            const out = {
                id: 1,
                username: "ExampleUser",
            } as UserData;
            if (Math.random() < 0.5) out.picture = Picture;
            return {data: out};
        }, //todo use actual endpoint
        // UpdateDataMy: async (args: UserData) => {
        //     console.log("API.Users.UpdateDataMy", args);
        // }, // todo use actual endpoint
        GetData: async (id: number) => {
            const out = {
                username: "ID:"+id,
            } as UserData;
            if (Math.random() < 0.5) out.picture = Picture;
            return {data: out};
        }, // todo use actual endpoint
        GetSettingsMy: async () => {
            const add = await API.Users.GetDataMy();
            const out = {
                ...add.data,
                interval: 60*60*24,
                email: "a@a.a",
            } as UserSettings;
            return {data: out};
        },
        UpdateSettingsMy: async (args: UserSettings) => {
            console.log("API.Users.UpdateSettingsMy", args);
        },
        // GetPicture: async (id: number) => ({data: Picture}) // todo use actual endpoint
        Search: () => {}, // todo
        Delete: () => {}, // todo
    }
};
