import * as d3 from "d3";
import type { ForecastDiff } from '../../../types';
import { createSvg, getContainerDimensions, throwAwayOldGraph } from '../../d3/graphUtils';
import { createXScale, createYScale, getXDomain, getYDomain } from './domains';
import type { D3Selection } from '../../../types.d3';

export const DIFF_GRAPH_ROOT = 'diff-graph'
const GRAPH_ROOT_SELECTOR = '.' + DIFF_GRAPH_ROOT

function drawLine(svg: D3Selection<SVGElement>, diff: ForecastDiff, line, x, y) {


  svg.append("path")
    .datum(diff.diffs)
    .attr("fill", "none")
    .attr("stroke", "steelblue")
    .attr("stroke-width", 1.5)
    .attr("d", line)
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
  const y = createYScale(bottomEdge, getYDomain(diffs))
  const svg = createSvg(GRAPH_ROOT_SELECTOR, width, height, margin)

  const line = d3.line().x(
    function(d,i) { return x(i)}
  ).y(
    function(d) { return y(d as any) }
  )
  diffs.forEach(d => {
    drawLine(svg, d, line, x, y)
  })
}
