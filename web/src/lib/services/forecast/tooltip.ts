import * as d3 from 'd3';
import type { D3Selection, Forecast } from '../../types';
import { formatDay, formatTime } from '../../utils/date';
import { getPeriodStart } from './domains';
import { isSmallViewport } from '../../utils/dom';

const MARGIN = 5

export function createTooltip(svg: D3Selection<SVGElement>, y: number) {
  return svg
    .append("text")
    .attr("class", "tooltip")
    .attr('x', MARGIN)
    .attr('y', y + MARGIN)
    .attr('font-size', isSmallViewport() ? '7px' : '15px')
    .attr('fill', '#ffffff')
}

export function getMouseOverHandler(tooltip: D3Selection<SVGTextElement>) {
  return function () {
    const {value, actual, periodEnd, actualCount}: Forecast = d3.select(this).datum() as Forecast;
    const dateTime = `${formatDay(periodEnd)} ${formatTime(getPeriodStart(periodEnd).toISOString())} - ${formatTime(periodEnd)}`
    const values = `Forecast ${value.toFixed(0)}W | Actual: ${actual === null ? 'No data' : actual.toFixed(0) + 'W'}`
    const count = `Count: ${actualCount}`
    tooltip.text(`${dateTime} | ${values} | ${count}`)
  }
}

export function getMouseLeaveHandler(tooltip: D3Selection<SVGTextElement>) {
  return () => tooltip.text('')
}