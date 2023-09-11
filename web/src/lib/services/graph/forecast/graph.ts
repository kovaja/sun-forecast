import type { Forecast } from '../../../types';
import * as d3 from "d3";
import { createXScale, createYScale, getXDomain, getYDomain } from './domains';
import { createTooltip } from './tooltip';
import { appendText } from './text';
import { formatTime } from '../../../utils/date';
import { getAnimateColumns, getAppendColumn } from './column';
import { addAttributes } from '../d3/utils';
import { createSvg, getContainerDimensions, throwAwayOldGraph } from '../d3/graph';
import { appendXAxis, appendXGrid, appendYAxis, appendYGrid, createXGrid, createYGrid } from '../d3/axis';

export const GRAPH_ROOT = 'graph-root';
const GRAPH_ROOT_SELECTOR = '.' + GRAPH_ROOT
const MAX_COLUMN_PADDING = 15
let COLUMN_PADDING = 15

function createColPadScale() {
  return d3.scaleLinear()
    .domain([10, 30])
    .range([MAX_COLUMN_PADDING, 3])
    .clamp(true);
}


function appendCurrentTimeIndicator(svg, x, bottomEdge) {
  const now = new Date()
  const xCoor = x(now)
  const yCoor = bottomEdge - 50

  addAttributes(svg.append("line"), {
    x1: xCoor,
    y1: yCoor,
    x2: xCoor,
    y2: bottomEdge + 10,
    'class': 'current-time-line'
  }).style("stroke-width", 3);


    addAttributes(svg.append("text"), {
      x: xCoor - 2,
      y: yCoor - 10,
      'font-size': 13,
      'class': 'current-time-text'
    })
    .text(formatTime(now.toISOString()))

}

export function plotGraph(data: Forecast[]) {
  const graphContainer = document.querySelector(GRAPH_ROOT_SELECTOR)

  if (!graphContainer) {
    return;
  }

  throwAwayOldGraph(graphContainer)

  const { width, height, rightEdge, bottomEdge, margin} = getContainerDimensions(graphContainer);

  const colPadScale = createColPadScale()
  COLUMN_PADDING = colPadScale(data.length)

  const x = createXScale(rightEdge, getXDomain(data))
  const y = createYScale(bottomEdge, getYDomain(data))
  const svg = createSvg(GRAPH_ROOT_SELECTOR, width, height, margin)

  const xAxisGrid = createXGrid(x, bottomEdge, 10)
  const yAxisGrid = createYGrid(y, rightEdge, 5)

  appendXGrid(svg, xAxisGrid, bottomEdge)
  appendYGrid(svg, yAxisGrid)

  const tooltip = createTooltip(svg, margin.top)
  const appendColumnsFn = getAppendColumn({
    elements: { svg, tooltip },
    dimensions: { rightEdge, bottomEdge, columnPadding: COLUMN_PADDING },
    scales: { x, y },
    data
  })
  const animateColumnsFn = getAnimateColumns({
    elements: { svg, tooltip },
    dimensions: { rightEdge, bottomEdge, columnPadding: COLUMN_PADDING },
    scales: { x, y },
  })
  appendColumnsFn('value')
  appendColumnsFn('actual')
  animateColumnsFn('value')
  animateColumnsFn('actual')

  appendXAxis(svg, bottomEdge, x)
  appendYAxis(svg, y)
  appendText(svg, data, x, y, COLUMN_PADDING)
  appendCurrentTimeIndicator(svg, x, bottomEdge)
}