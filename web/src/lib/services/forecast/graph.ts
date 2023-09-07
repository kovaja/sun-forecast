import type { Forecast } from '../../types';
import * as d3 from "d3";
import { createXScale, getXDomain, createYScale, getYDomain, getPeriodStart } from './domains';
import { createTooltip, getMouseLeaveHandler, getMouseOverHandler } from './tooltip';
import {
  ACTUAL_BAR_FILL,
  ACTUAL_BAR_STROKE,
  FORECAST_BAR_FILL,
  FORECAST_BAR_STROKE
} from './constants';
import { appendText } from './text';
import { formatTime } from '../../utils/date';

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

function getContainerDimensions(): { width: number, height: number} {
  const rect =  document.querySelector(GRAPH_ROOT_SELECTOR).getBoundingClientRect()
  // 480 is breakpoint for tabs disappearing
  const headerSpace = rect.width < 480 ? 70 : 120
  return {
    width: rect.width * 0.95,
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

function appendColumns(
  svg,
  rightEdge: number,
  bottomEdge: number,
  data: Forecast[],
  x,
  y,
  tooltip,
  property: 'value' | 'actual',
  fill: string,
  stroke: string,
) {

  svg.selectAll('bar-' + property)
    .data(data)
    .enter()
    .append("rect")
    .attr("x", function (d) {
      return COLUMN_PADDING / 2 + x(getPeriodStart(d.periodEnd));
    })
    .attr("y", function (d) {
      return y(d[property] ?? 0);
    })
    .attr("width", (rightEdge / data.length) - COLUMN_PADDING)
    .attr("height", function (d) {
      const val = d[property] ?? 0
      return val === 0 ? 0 : bottomEdge - y(val)
    })
    .attr("fill", fill)
    .attr("stroke", stroke)
    .attr("rx", "5")
    .on("mouseover", getMouseOverHandler(tooltip))
    // .on("mousemove", mousemove)
    .on("mouseleave", getMouseLeaveHandler(tooltip))
}

function appendCurrentTimeIndicator(svg, x, bottomEdge) {
  const now = new Date()
  const xCoor = x(now)
  const yCoor = 50

  svg.append("line")
    .attr("x1", xCoor)
    .attr("y1", yCoor)
    .attr("x2", xCoor)
    .attr("y2", bottomEdge)
    .attr("class", "current-time-line")
    .style("stroke-width", 1);

  svg.append("text")
    .attr('x', xCoor - 2)
    .attr('y', yCoor - 10)
    .attr('font-size', 13)
    .attr('class', 'current-time-text')
    .text(formatTime(now.toISOString()))

}

function throwAwayOldGraph() {
  const graphRoot = document.querySelector(GRAPH_ROOT_SELECTOR)
  if (graphRoot) {
    graphRoot.innerHTML = ''
  }
}

export function plotGraph(data: Forecast[]) {
  throwAwayOldGraph()

  const margin = {top: 10, left: 35, right: 5, bottom: 40};
  const { width, height } = getContainerDimensions();
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
  appendColumns(svg, rightEdge, bottomEdge, data, x, y, tooltip, 'value', FORECAST_BAR_FILL, FORECAST_BAR_STROKE)
  appendColumns(svg, rightEdge, bottomEdge, data, x, y, tooltip, 'actual', ACTUAL_BAR_FILL, ACTUAL_BAR_STROKE)
  appendXAxis(svg, bottomEdge, x)
  appendYAxis(svg, y)
  appendText(svg, data, x, y, COLUMN_PADDING)
  appendCurrentTimeIndicator(svg, x, bottomEdge)
}