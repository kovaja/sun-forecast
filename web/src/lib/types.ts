export interface ApiResponse<T> {
    date: string;
    num: number;
    data: T | null
}

export interface AppEvent {
    message: string;
    timestamp: string;
}