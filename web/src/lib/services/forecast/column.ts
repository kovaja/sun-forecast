import type { D3LinearScale, D3Selection, D3TimeScale } from '../../types.d3';
import type { Forecast } from '../../types';
import { getPeriodStart } from './domains';
import { getMouseLeaveHandler, getMouseOverHandler } from './tooltip';
import { ACTUAL_BAR_FILL, ACTUAL_BAR_STROKE, FORECAST_BAR_FILL, FORECAST_BAR_STROKE } from './constants';
import type { BaseType } from 'd3';
import { addAttributes } from './utils.d3';

type ColumnProperty = 'value' | 'actual'
const FILL_MAP: Record<ColumnProperty, string> = {
  value: FORECAST_BAR_FILL,
  actual: ACTUAL_BAR_FILL
}
const STROKE_MAP: Record<ColumnProperty, string> = {
  value: FORECAST_BAR_STROKE,
  actual: ACTUAL_BAR_STROKE
}

const SELECTOR_MAP: Record<ColumnProperty, string> = {
  value: 'bar-value',
  actual: 'bar-actual'
}

interface ColumnParams {
  elements: {
    svg: D3Selection<SVGElement>;
    tooltip: D3Selection<SVGTextElement>,
  }
  dimensions: {
    rightEdge: number,
    bottomEdge: number,
    columnPadding: number,
  }
  data: Forecast[],
  scales: {
    x: D3TimeScale,
    y: D3LinearScale
  }
}



export function getAppendColumn({
 elements,
 data,
 dimensions,
 scales,
}: ColumnParams): (property: ColumnProperty) => void {
  return (property: ColumnProperty) => {
    const {svg, tooltip} = elements
    const {rightEdge, bottomEdge, columnPadding} = dimensions
    const {x, y} = scales
    const fill = FILL_MAP[property]
    const stroke = STROKE_MAP[property]

    function getXCoord(d: Forecast): number {
      return columnPadding / 2 + x(getPeriodStart(d.periodEnd));
    }

    function getYCoord(d: Forecast): number {
      return y(0);
    }

    function getHeight(_: Forecast): number {
      return 0
    }

    const width = (rightEdge / data.length) - columnPadding

    const selection = svg.selectAll('.' + SELECTOR_MAP[property])
      .data(data)
      .enter()
      .append("rect")

    addAttributes(selection, {
      'class': SELECTOR_MAP[property],
      x: getXCoord,
      y: getYCoord,
      width: width,
      height: getHeight,
      fill: fill,
      stroke: stroke,
      rx: "5",
    })
      .on("mouseover", getMouseOverHandler(tooltip))
      .on("mouseleave", getMouseLeaveHandler(tooltip))
  }
}

export function getAnimateColumns(params: Pick<ColumnParams, 'elements' | 'scales' | 'dimensions'>): (property: ColumnProperty) => void {
  const { svg } = params.elements
  const { y } = params.scales
  const { bottomEdge } = params.dimensions

  return (property: ColumnProperty) => {
    function getYCoord(d: Forecast): number {
      return y(d[property] ?? 0);
    }

    function getHeight(d: Forecast): number {
      const val = d[property] ?? 0
      return val === 0 ? 0 : bottomEdge - y(val)
    }

    svg.selectAll('.' + SELECTOR_MAP[property])
      .transition()
      .duration(800)
      .attr("y", getYCoord)
      .attr("height", getHeight)
      .delay((_, i) => i * 50)
  }
}
