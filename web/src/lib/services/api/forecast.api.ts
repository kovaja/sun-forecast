import type { Forecast } from '../../types';
import { fetchJsonData } from './base';

const MS_3_HOURS = 3* 60*60*1000
const MS_6_HOURS = 6* 60*60*1000

export function fetchForecast(): Promise<Forecast[] | null> {
    const nowTs = new Date().getTime()
    const params = {
        from: new Date(nowTs - MS_3_HOURS).toISOString(),
        to: new Date(nowTs+MS_6_HOURS).toISOString()
    }
    return fetchJsonData<Forecast[]>('forecast', params)
}