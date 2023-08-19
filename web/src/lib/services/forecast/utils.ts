import type { Forecast } from '../../types';

export function forecastSuccessRate(forecast: Forecast): number | null {
    if (forecast.value === 0 || forecast.actual === null) {
        return null
    }
    return -1 * (100 - (forecast.actual / forecast.value) * 100)
}