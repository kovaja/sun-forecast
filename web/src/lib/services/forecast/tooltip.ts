import * as d3 from 'd3';
import type { Forecast } from '../../types';
import { formatDay, formatTime } from '../../utils/date';
import { getPeriodStart } from './domains';
import type { Selection } from 'd3';

const MARGIN = 5

export function createTooltip(svg: Selection<SVGElement, any, any, any>, y: number) {
  return svg
    .append("text")
    // .style("opacity", 0)
    .attr("class", "tooltip")
    .attr('x', MARGIN)
    .attr('y', y + MARGIN)
    .attr('width', '100%')
    .attr('height', 25)
    .attr('fill', '#ffffff')
}

// Three function that change the tooltip when user hover / move / leave a cell
export function getMouseOverHandler(tooltip) {
  return function () {
    const {value, actual, periodEnd}: Forecast = d3.select(this).datum() as Forecast;
    const dateTime = `${formatDay(periodEnd)} ${formatTime(getPeriodStart(periodEnd).toISOString())} - ${formatTime(periodEnd)}`
    const values = `Forecast ${value.toFixed(0)}W | Actual: ${actual === null ? 'No data' : actual.toFixed(0) + 'W'}`
    tooltip
      .text(`${dateTime}: ${values}`)
    // .style("opacity", 1)
  }
}

export function getMouseLeaveHandler(tooltip) {
  return () => tooltip.text('')
}