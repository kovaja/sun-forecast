import * as d3 from "d3";
import type { ForecastDiff } from '../../../types';

export function getYDomain(diffs: ForecastDiff[]) {
  const flatDiffs = diffs.map(d => d.diffs).flat()
  return [d3.min(flatDiffs), d3.max(flatDiffs)]
}
export function getXDomain(diffs: ForecastDiff[]) {
  return [0, diffs[0].diffs.length -1]
}

export function createXScale(rightEdgeX: number, domain: number[]) {
  return d3.scaleLinear()
    .domain(domain)
    .range([0, rightEdgeX])
}

export function createYScale(bottomEdgeY: number, domain: number[]) {
  return d3.scaleLinear()
    .domain(domain)
    .range([bottomEdgeY, 0]);
}