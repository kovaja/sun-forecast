import type { AppEvent } from '../../types';
import { fetchJsonData } from './base';
import { AppEventType } from '../../types';

export async function fetchEvents(type: AppEventType): Promise<AppEvent[] | null> {
  return fetchJsonData<AppEvent[]>('event', { type })
}