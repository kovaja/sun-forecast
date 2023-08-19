import type { ApiResponse, AppEvent } from '../../types';
import { logger } from '../logger';

const API_BASE = '/api'

export async function fetchEvents(): Promise<AppEvent[] | null> {
    const API_PATH = '/event/'
    const url = API_BASE + API_PATH

    logger.log('Fetch', url)

    try {
        const response = await fetch(url)
        const result: ApiResponse<AppEvent[]> = await response.json()

        return result.data
    } catch (e) {
        logger.log('Failed to read events', e)
        return null
    }
}