export interface DataResponse<T> {
  date: string;
  num: number;
  data: T | null
}

export interface ErrorResponse {
  error: string;
}

export type ApiResponse<T> = DataResponse<T> | ErrorResponse

export enum AppEventType {
  ForecastConsumed,
  ForecastUpdated,
  AppError,
}
export interface AppEvent {
  message: string;
  timestamp: string;
  type: AppEventType;
}

export interface Forecast {
  id: number;
  periodEnd: string;
  value: number;
  actual: number;
  actualCount: number;
}

export interface ForecastResponse {
  forecasts: Forecast[] | null;
  from: string;
  to: string;
}

export interface AppLink {
  route: string;
  name: string;
}