export function formatDay(inputDate: string): string {
  const date = new Date(inputDate);

  const day = String(date.getDate()).padStart(2, '0');
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const year = date.getFullYear();

  return `${day}.${month}.${year}`
}

export function formatTime(inputDate: string): string {
  const date = new Date(inputDate);

  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');

  return `${hours}:${minutes}`
}

export function formatDate(inputDate: string): string {
  return `${formatDay(inputDate)} ${formatTime(inputDate)}`;
}