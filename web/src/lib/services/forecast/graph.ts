import type { Forecast } from '../../types';
import * as d3 from "d3";
import { createXScale, getXDomain, createYScale, getYDomain } from './domains';
import { createTooltip } from './tooltip';
import { appendText } from './text';
import { formatTime } from '../../utils/date';
import { isSmallViewport } from '../../utils/dom';
import { getAnimateColumns, getAppendColumn } from './column';
import { addAttributes } from './utils.d3';

export const GRAPH_ROOT = 'graph-root';
const GRAPH_ROOT_SELECTOR = '.' + GRAPH_ROOT
const MAX_COLUMN_PADDING = 15
let COLUMN_PADDING = 15

interface Margin {
  top: number;
  left: number;
  bottom: number;
  right: number;
}

function createColPadScale() {
  return d3.scaleLinear()
    .domain([10, 30])
    .range([MAX_COLUMN_PADDING, 3])
    .clamp(true);
}

function getContainerDimensions(container: Element): { width: number, height: number } {
  const rect = container.getBoundingClientRect()
  const headerSpace = isSmallViewport() ? 80 : 120
  return {
    width: rect.width,
    height: window.innerHeight - headerSpace + 10
  }
}

function createSvg(width: number, height: number, margin: Margin) {
  return d3.select(GRAPH_ROOT_SELECTOR)
    .append("svg")
    .attr("width", width)
    .attr("height", height)
    .append("g")
    .attr(
      "transform",
      "translate(" + margin.left + "," + margin.top + ")"
    );
}

function createXGrid(x, size, numberOfTicks) {
  return d3.axisBottom(x).tickSize(-size).tickFormat(() => '').ticks(numberOfTicks);
}

function createYGrid(y, size, numberOfTicks) {
  return d3.axisLeft(y).tickSize(-size).tickFormat(() => '').ticks(numberOfTicks);
}


function appendXGrid(svg, grid, height) {
  svg.append('g')
    .attr('class', 'x axis-grid')
    .attr('transform', 'translate(0,' + height + ')')
    .call(grid)
}

function appendYGrid(svg, grid) {
  svg.append('g')
    .attr('class', 'y axis-grid')
    .call(grid);
}

function appendXAxis(svg, bottomEdge: number, x) {
  svg.append("g")
    .attr("transform", "translate(0," + bottomEdge + ")")
    .call(d3.axisBottom(x))
    .selectAll("text")
    .attr("transform", "translate(-10,0)rotate(-45)")
    .style("text-anchor", "end");
}

function appendYAxis(svg, y) {
  svg.append("g")
    .call(d3.axisLeft(y));
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

function throwAwayOldGraph(graphRoot: Element) {
  graphRoot.innerHTML = ''
}

export function plotGraph(data: Forecast[]) {
  const graphContainer = document.querySelector(GRAPH_ROOT_SELECTOR)

  if (!graphContainer) {
    return;
  }

  throwAwayOldGraph(graphContainer)

  const margin = {top: 10, left: 35, right: 5, bottom: 40};
  const {width, height} = getContainerDimensions(graphContainer);
  const rightEdge = width - margin.left - margin.right
  const bottomEdge = height - margin.top - margin.bottom

  const colPadScale = createColPadScale()
  COLUMN_PADDING = colPadScale(data.length)

  const x = createXScale(rightEdge, getXDomain(data))
  const y = createYScale(bottomEdge, getYDomain(data))
  const svg = createSvg(width, height, margin)

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