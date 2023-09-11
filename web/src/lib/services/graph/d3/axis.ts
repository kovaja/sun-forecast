import * as d3 from 'd3';

export function createXGrid(x, size, numberOfTicks) {
  return d3.axisBottom(x).tickSize(-size).tickFormat(() => '').ticks(numberOfTicks);
}

export function createYGrid(y, size, numberOfTicks) {
  return d3.axisLeft(y).tickSize(-size).tickFormat(() => '').ticks(numberOfTicks);
}

export function appendXGrid(svg, grid, height) {
  svg.append('g')
    .attr('class', 'x axis-grid')
    .attr('transform', 'translate(0,' + height + ')')
    .call(grid)
}

export function appendYGrid(svg, grid) {
  svg.append('g')
    .attr('class', 'y axis-grid')
    .call(grid);
}

export function appendXAxis(svg, bottomEdge: number, x, tickFormatFn = (d) => d) {
  svg.append("g")
    .attr("transform", "translate(0," + bottomEdge + ")")
    .call(d3.axisBottom(x).tickFormat(tickFormatFn))
    .selectAll("text")
    .attr("transform", "translate(-10,0)rotate(-45)")
    .style("text-anchor", "end");
}

export function appendYAxis(svg, y) {
  svg.append("g")
    .call(d3.axisLeft(y));
}