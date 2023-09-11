import type { D3Selection } from '../../../types.d3';
import type { ForecastDiff } from '../../../types';
import * as d3 from 'd3';

interface DrawLineParms {
  svg: D3Selection<SVGElement>;
  diff: ForecastDiff;
  line: any; // d3.Line<[number, number]>;
  stroke: string;
  strokeWidth: number;
}

export function drawLine({
  svg,
  diff,
  line,
  stroke,
  strokeWidth
}: DrawLineParms) {
  svg.append("path")
    .datum(diff.diffs)
    .attr("fill", "none")
    .attr("stroke", stroke)
    .attr("stroke-width", strokeWidth)
    .attr("d", line)
    .on('mouseover', function () {
      d3.select(this).attr('stroke-width', strokeWidth * 4)
    })
    .on('mouseout', function (){
      d3.select(this).attr('stroke-width', strokeWidth)
    })
}