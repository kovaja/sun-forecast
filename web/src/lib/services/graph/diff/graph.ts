import * as d3 from "d3";
import type { ForecastDiff } from '../../../types';
import { createSvg, getContainerDimensions, throwAwayOldGraph } from '../d3/graph';
import { createXScale, createYScale, getXDomain, getYDomain } from './domains';
import type { D3Selection } from '../../../types.d3';
import { appendXAxis, appendXGrid, appendYAxis, appendYGrid, createXGrid, createYGrid } from '../d3/axis';
import { AVG_LINE_STROKE, DIFF_LINE_STROKE } from './constants';
import { drawLine } from './line';
import { computeAverageSeries } from './averageLine';

export const DIFF_GRAPH_ROOT = 'diff-graph'
const GRAPH_ROOT_SELECTOR = '.' + DIFF_GRAPH_ROOT



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
      strokeWidth: 0.5
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

  appendXAxis(svg, bottomEdge, x)
  appendYAxis(svg, y)
}
