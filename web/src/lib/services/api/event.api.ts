import type { AppEvent } from '../../types';
import { fetchJsonData } from './base';

export async function fetchEvents(): Promise<AppEvent[] | null> {
    return fetchJsonData<AppEvent[]>('event')
}