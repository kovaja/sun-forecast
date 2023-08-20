import * as d3 from 'd3';
import type { Forecast } from '../../types';
import type { FadeParams } from 'svelte/transition';
import { formatDate } from '../../utils/date';

export function createTooltip(selector: string) {
  return d3.select(selector)
    .append("div")
    .style("opacity", 0)
    .attr("class", "tooltip")
    .style("background-color", "#395B64")
    .style("position", "absolute")
    .style("border", "solid")
    .style("border-width", "1px")
    .style("border-radius", "5px")
    .style("padding", "5px")
    .style("top", 200)
    .style('right', 0)
    .style('width', '30%');
}

// Three function that change the tooltip when user hover / move / leave a cell
export function getMouseOverHandler(tooltip) {
  return function () {
    const {value, actual, periodEnd }: Forecast = d3.select(this).datum() as Forecast;
    const html = `
     <ul>
        <li>Time: ${formatDate(periodEnd)}</li>
       <li>Forcast: ${value.toFixed(0)}</li>
       <li>Actual: ${actual === null ? 'No data' : actual.toFixed(0)}</li>
     </ul>
    `
    tooltip
      .html(html)
      .style("opacity", 1)
  }
}

// const mousemove = function(d) {
//     tooltip
//         .style("left", (d3.mouse(this)[0]+90) + "px") // It is important to put the +90: other wise the tooltip is exactly where the point is an it creates a weird effect
//         .style("top", (d3.mouse(this)[1]) + "px")
// }
export function getMouseLeaveHandler(tooltip) {
  return () => {
    tooltip
      .style("opacity", 0)
  }
}