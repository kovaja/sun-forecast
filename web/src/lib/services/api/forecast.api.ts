import type { ForecastDiff, ForecastParams, ForecastResponse } from '../../types';
import { fetchJsonData } from './base';


const MS_DAYS = (x: number) => x * 24 * 60 * 60 * 1000

export async function fetchForecast(params: ForecastParams): Promise<ForecastResponse | null> {
  return await fetchJsonData<ForecastResponse>('forecast', params)
}

export async function fetchDiffs(days: number): Promise<ForecastDiff[]> {
  const params = {
    from: new Date(new Date().getTime() - MS_DAYS(days)).toISOString(),
    to: new Date().toISOString()
  }
  return await fetchJsonData<ForecastDiff[]>('forecast/diff', params)
}