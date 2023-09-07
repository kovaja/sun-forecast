import { getPeriodStart } from './domains';
import * as d3 from 'd3';
import { forecastSuccessRate } from './utils';
import { SUCCESSRATE_GREEN, SUCCESSRATE_RED } from './constants';

export function appendText(svg, data, x, y, columnPadding) {
  svg.selectAll("successRate")
    .data(data)
    .enter()
    .append("text")
    .attr("x", function (d) {
      return columnPadding / 2 + x(getPeriodStart(d.periodEnd)) + 5;
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
        return successRate > 0 ?  SUCCESSRATE_GREEN : SUCCESSRATE_RED
      }
      return ''
    })
    .attr("font-size", "12px")
    .style("opacity", 0)
    .transition()
    .duration(1600)
    .style("opacity", 1)
    .delay((_, i) => i * 50+5)
}
