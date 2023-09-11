import * as d3 from "d3";
import type { ForecastDiff } from '../../../types';
import { createSvg, getContainerDimensions, throwAwayOldGraph } from '../d3/graph';
import { createXScale, createYScale, getXDomain, getYDomain } from './domains';
import { appendXAxis, appendXGrid, appendYAxis, appendYGrid, createXGrid, createYGrid } from '../d3/axis';
import { AVG_LINE_STROKE, DIFF_LINE_STROKE, MEDIAN_LINE_STROKE } from './constants';
import { drawLine } from './line';
import { computeAverageSeries, computeMedianSeries } from './aggregateLine';

export const DIFF_GRAPH_ROOT = 'diff-graph'
const GRAPH_ROOT_SELECTOR = '.' + DIFF_GRAPH_ROOT

/**
 * Converts 0,1,2... to periodEnd (00:00, 00:30, 01:00, ...)
 * @param d
 */
function formatTick(d: any) {
  const hrs = d / 2;
  const hasHalf = hrs*10 % 10 !== 0
  const fullHr = Math.floor(hrs)

  return `${String(fullHr).padStart(2, '0')}:${hasHalf ? '30' : '00'}`
}

export function plotGraph(diffs: ForecastDiff[]) {
  const graphContainer = document.querySelector(GRAPH_ROOT_SELECTOR)
  if (!graphContainer) {
    return;
  }

  throwAwayOldGraph(graphContainer)

  const {
    width,
    height,
    rightEdge,
    bottomEdge,
    margin
  } = getContainerDimensions(graphContainer);

  const x = createXScale(rightEdge, getXDomain(diffs))
  const y = createYScale(bottomEdge, getYDomain(diffs ))
  const svg = createSvg(GRAPH_ROOT_SELECTOR, width, height, margin)

  const xAxisGrid = createXGrid(x, bottomEdge, 10)
  const yAxisGrid = createYGrid(y, rightEdge, 5)

  appendXGrid(svg, xAxisGrid, bottomEdge)
  appendYGrid(svg, yAxisGrid)

  const line = d3.line().x(
    function(d,i) { return x(i)}
  ).y(
    function(d) { return y(d as any) }
  )
  diffs.forEach(d => {
    drawLine({
      svg,
      diff: d,
      line,
      stroke: DIFF_LINE_STROKE,
      strokeWidth: 0.3
    })
  })

  const avgSeries = computeAverageSeries(diffs)
  drawLine({
    svg,
    diff: { diffs: avgSeries, date: 'Average' },
    line,
    stroke: AVG_LINE_STROKE,
    strokeWidth: 2
  })

  const medSeries = computeMedianSeries(diffs)
  drawLine({
    svg,
    diff: { diffs: medSeries, date: 'Median' },
    line,
    stroke: MEDIAN_LINE_STROKE  ,
    strokeWidth: 2
  })

  appendXAxis(svg, bottomEdge, x, formatTick)
  appendYAxis(svg, y)
}
