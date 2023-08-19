function log(...args: any[]) {
    console.log('[SunForecast]', ...args)
}

export const logger = { log }