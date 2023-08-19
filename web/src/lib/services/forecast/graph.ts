import type { Forecast } from '../../types';
import * as d3 from "d3";
import { createXScale, getXDomain, createYScale, getYDomain, getPeriodStart } from './domains';
import { createTooltip, getMouseLeaveHandler, getMouseOverHandler } from './tooltip';
import { forecastSuccessRate } from './utils';

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

function getContainerWidth(): number {
    return document.querySelector(GRAPH_ROOT_SELECTOR).getBoundingClientRect().width
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
        .on("mouseover", getMouseOverHandler(tooltip))
        // .on("mousemove", mousemove)
        .on("mouseleave", getMouseLeaveHandler(tooltip))
}

function appendText(svg, data, x, y) {
    svg.selectAll("successRate")
        .data(data)
        .enter()
        .append("text")
        .attr("x", function (d) {
            return COLUMN_PADDING / 2 + x(getPeriodStart(d.periodEnd));
        })
        .attr("y", function (d) {
            return d3.min([y(d.actual ?? 0), y(d.value)]) - 10;
        })
        .text(d => {
            const successRate = forecastSuccessRate(d)
            if (typeof successRate === 'number') {
                return successRate.toFixed(0) + '%'
            }
            return ''
        })
        .attr("fill", d => {
            const successRate = forecastSuccessRate(d)
            if (typeof successRate === 'number') {
                return successRate > 0 ? 'green' : 'red'
            }
            return ''
        })
        .attr("font-size", "12px")
}

function throwAwayOldGraph() {
    document.querySelector(GRAPH_ROOT_SELECTOR).innerHTML = ''
}

export function plotGraph(data: Forecast[]) {
    throwAwayOldGraph()

    const margin = {top: 10, left: 35, right: 5, bottom: 70};
    const width = getContainerWidth();
    const height = 600 - margin.top - margin.bottom;
    const rightEdge = width - margin.left - margin.right
    const bottomEdge = height - 50

    const colPadScale = createColPadScale()
    COLUMN_PADDING = colPadScale(data.length)
    console.log('colpad', COLUMN_PADDING)

    const tooltip = createTooltip(GRAPH_ROOT_SELECTOR)
    const x = createXScale(rightEdge, getXDomain(data))
    const y = createYScale(bottomEdge, getYDomain(data))
    const svg = createSvg(width, height, margin)

    appendXAxis(svg, bottomEdge, x)
    appendYAxis(svg, y)
    appendColumns(svg, rightEdge, bottomEdge, data, x, y, tooltip, 'value', "#69b3a2", "#69b3a2")
    appendColumns(svg, rightEdge, bottomEdge, data, x, y, tooltip, 'actual', "rgba(194,137,33,0.29)", "#c28921")
    appendText(svg, data, x, y)
}