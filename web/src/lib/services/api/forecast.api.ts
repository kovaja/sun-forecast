import type { Forecast, ForecastResponse } from '../../types';
import { fetchJsonData } from './base';

const MS_HOURS = (x) => x * 60 * 60 * 1000

export async function fetchForecast(windowSizeHrs: number, windowMiddle: Date): Promise<ForecastResponse | null> {
  const nowTs = windowMiddle.getTime()
  const offset = MS_HOURS(windowSizeHrs / 2)
  const params= {
    from: new Date(nowTs - offset).toISOString(),
    to: new Date(nowTs + offset).toISOString()
  }

  const forecasts: Forecast[] = await fetchJsonData<Forecast[]>('forecast', params)

  // temporary, TODO return from, to from backend
  return {
    data: forecasts,
    from: params.from,
    to: params.to
  }
}