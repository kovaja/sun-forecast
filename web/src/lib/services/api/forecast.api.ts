import type { Forecast } from '../../types';
import { fetchJsonData } from './base';

const MS_HOURS = (x) => x * 60*60*1000

export function fetchForecast(windowSizeHrs: number, windowMiddle: Date): Promise<Forecast[] | null> {
    const nowTs = windowMiddle.getTime()
    const offset = MS_HOURS(windowSizeHrs/2)
    const params = {
        from: new Date(nowTs - offset).toISOString(),
        to: new Date(nowTs + offset).toISOString()
    }
    return fetchJsonData<Forecast[]>('forecast', params)
}