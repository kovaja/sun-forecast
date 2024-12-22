import * as d3 from 'd3';

export interface Margin {
  top: number;
  left: number;
  bottom: number;
  right: number;
}

interface GraphDimension {
  width: number;
  height: number;
  rightEdge: number;
  bottomEdge: number;
  margin: Margin;
}

export function throwAwayOldGraph(graphRoot: Element) {
  graphRoot.innerHTML = ''
}

export function getContainerDimensions(container: Element, isGlobal: boolean): GraphDimension {
  const margin = {top: 10, left: 35, right: 5, bottom: 40};
  const rect = container.getBoundingClientRect()

  const width =  rect.width
  const height = (isGlobal ? (window.innerHeight - 80) : rect.height) + 10
  const rightEdge = width - margin.left - margin.right
  const bottomEdge = height - margin.top - margin.bottom


  return {
    width,
    height,
    rightEdge,
    bottomEdge,
    margin
  }
}

export function createSvg(selector: string, width: number, height: number, margin: Margin) {
  return d3.select(selector)
    .append("svg")
    .attr("width", width)
    .attr("height", height)
    .append("g")
    .attr(
      "transform",
      "translate(" + margin.left + "," + margin.top + ")"
    );
}