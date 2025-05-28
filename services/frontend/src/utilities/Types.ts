export interface Event {
    id?: number,
    name: string,
    description: string,
    category: string,
    participants?: number,
    max_participants: number,
    image_data?: string,
    city?: string,
    address?: string,
    latitude?: number,
    longitude?: number,
    start_time: string,
    end_time: string,
    created_by?: number,
};

export interface Registration {
    id: number,
    event_id: number,
    user_id: number,
    registration_time: string,
    // status: string,
    // updated_at: string,
    // comment?: string,
};

export interface Review {
    id: number,
    user_id: number,
    content: string,
    posted_at: string,
    updated_at: string,
};

export interface UserData {
    // id: number,
    username: string,
    picture?: string,
};

export interface UserSettings extends UserData {
    //interval: number, // не надо
    telegram?: string, // telegram id
    passhash?: string, // разрешить ли?
};

export interface UserCredentials {
    username: string,
    passhash: string,
};
