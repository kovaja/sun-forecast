import type { ForecastParams } from '../types';

export function formatDay(inputDate: string): string {
  const date = new Date(inputDate);

  const day = String(date.getDate()).padStart(2, '0');
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const year = date.getFullYear();

  return `${day}.${month}.${year}`
}

export function formatTime(inputDate: string, includeSeconds = false): string {
  const date = new Date(inputDate);

  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');

  return `${hours}:${minutes}${includeSeconds ? ':' + seconds : ''}`
}

export function formatDate(inputDate: string, includeSeconds = false): string {
  return `${formatDay(inputDate)} ${formatTime(inputDate, includeSeconds)}`;
}

const MS_HOURS = (x: number) => x * 60 * 60 * 1000

export function mapWindowToForecastParams(windowSizeHrs: number, windowMiddle: Date): ForecastParams {
  const nowTs = windowMiddle.getTime()
  const offset = MS_HOURS(windowSizeHrs / 2)
  return {
    from: new Date(nowTs - offset).toISOString(),
    to: new Date(nowTs + offset).toISOString()
  }
}

export function getYesterdayStartMs(): number {
  const today = new Date();
  const yesterday = new Date(today.getFullYear(), today.getMonth(), today.getDate() - 1);
  return yesterday.getTime()
}

export function getTodayStartMs(): number {
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  return today.getTime()
}