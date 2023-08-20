export interface DataResponse<T> {
  date: string;
  num: number;
  data: T | null
}

export interface ErrorResponse {
  error: string;
}

export type ApiResponse<T> = DataResponse<T> | ErrorResponse

export interface AppEvent {
  message: string;
  timestamp: string;
}

export interface Forecast {
  id: number;
  periodEnd: string;
  value: number;
  actual: number;
}

export interface AppLink {
  route: string;
  name: string;
}