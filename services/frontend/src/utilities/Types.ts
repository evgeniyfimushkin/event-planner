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
};

export interface Registration {
    id: number,
    event_id: number,
    user_id: number,
    registration_time: string,
    status: string,
    updated_at: string,
    comment?: string,
};

export interface Review {
    id: number,
    user_id: number,
    content: string,
    posted_at: string,
    updated_at: string,
}

export interface UserData {
    id: number,
    username: string,
    picture?: string,
}
