import { logger } from '../logger';
import type { ApiResponse, DataResponse } from '../../types';

const API_BASE = '/api'

type SupportedPath = 'event' | 'forecast'

function getUrl(path: SupportedPath, query?: string): string {
  const url = `${API_BASE}/${path}/`
  if (query) {
    return url + '?' + query
  }
  return url
}

function isDataResponse<T>(x: unknown): x is DataResponse<T> {
  return !!x && x.hasOwnProperty('date')
}

export async function fetchJsonData<T>(path: SupportedPath, queryParams?: Record<string, string>): Promise<T | null> {
  const query = queryParams ? new URLSearchParams(queryParams).toString() : ''
  const url = getUrl(path, query)
  logger.log('Fetch', url)

  try {
    const response = await fetch(url)
    const result: ApiResponse<T> = await response.json()

    if (response.status !== 200 && !isDataResponse(result)) {
      logger.log('Api error:', result.error)
      return null
    }

    if (isDataResponse(result)) {
      logger.log(`${path} returned ${result.num} records`)
      return result.data
    }

    logger.log('Fetch result is not api result')
    return null
  } catch (e) {
    logger.log('Failed to fetch', url, e)
    return null
  }
}