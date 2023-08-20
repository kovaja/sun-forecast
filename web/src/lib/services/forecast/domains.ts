import type { Forecast } from '../../types';
import * as d3 from 'd3';

const MS_30_MINUTES = 30 * 60 * 1000
const DOMAIN_BLOAT = 1.20 // multiplier for Y domain to get more space above columns

export function getPeriodStart(periodEnd: string): Date {
  const endMs = new Date(periodEnd).getTime()
  return new Date(endMs - MS_30_MINUTES)
}

function getLastPeriod(data: Forecast[]): string {
  return data[data.length - 1].periodEnd
}

function getFirstPeriod(data: Forecast[]): string {
  return data[0].periodEnd
}

export function getXDomain(data: Forecast[]): Date[] {
  const firstPeriod = getFirstPeriod(data)
  const lastPeriod = getLastPeriod(data)
  return [
    getPeriodStart(firstPeriod),
    new Date(lastPeriod)
  ]
}

export function createXScale(rightEdgeX: number, domain: Date[]) {
  return d3.scaleTime()
    .domain(domain)
    .range([0, rightEdgeX])
}

export function getYDomain(data: Forecast[]): number[] {
  return [0, d3.max(data, d => d3.max([d.value, d.actual])) * DOMAIN_BLOAT]
}

export function createYScale(bottomEdgeY: number, domain: number[]) {
  return d3.scaleLinear()
    .domain(domain)
    .range([bottomEdgeY, 0]);
}