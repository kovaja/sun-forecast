import type { ForecastDiff, ForecastResponse } from '../../types';
import { fetchJsonData } from './base';

const MS_HOURS = (x) => x * 60 * 60 * 1000

export async function fetchForecast(windowSizeHrs: number, windowMiddle: Date): Promise<ForecastResponse | null> {
  const nowTs = windowMiddle.getTime()
  const offset = MS_HOURS(windowSizeHrs / 2)
  const params= {
    from: new Date(nowTs - offset).toISOString(),
    to: new Date(nowTs + offset).toISOString()
  }

  return await fetchJsonData<ForecastResponse>('forecast', params)
}

export async function fetchDiffs(): Promise<ForecastDiff[]> {
  const params = {
    from: new Date(new Date().getTime() - 6*24*60*60*1000).toISOString(),
    to: new Date().toISOString()
  }
  return await fetchJsonData<ForecastDiff[]>('forecast/diff', params)
}