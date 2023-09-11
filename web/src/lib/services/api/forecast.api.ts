import type { ForecastDiff, ForecastResponse } from '../../types';
import { fetchJsonData } from './base';

const MS_HOURS = (x: number) => x * 60 * 60 * 1000
const MS_DAYS = (x: number) => x * 24 * 60 * 60 * 1000

export async function fetchForecast(windowSizeHrs: number, windowMiddle: Date): Promise<ForecastResponse | null> {
  const nowTs = windowMiddle.getTime()
  const offset = MS_HOURS(windowSizeHrs / 2)
  const params= {
    from: new Date(nowTs - offset).toISOString(),
    to: new Date(nowTs + offset).toISOString()
  }

  return await fetchJsonData<ForecastResponse>('forecast', params)
}

export async function fetchDiffs(days: number): Promise<ForecastDiff[]> {
  const params = {
    from: new Date(new Date().getTime() - MS_DAYS(days)).toISOString(),
    to: new Date().toISOString()
  }
  return await fetchJsonData<ForecastDiff[]>('forecast/diff', params)
}