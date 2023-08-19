import * as d3 from 'd3';
import type { Forecast } from '../../types';
import type { FadeParams } from 'svelte/transition';

export function createTooltip(selector: string) {
    return d3.select(selector)
        .append("div")
        .style("opacity", 0)
        .attr("class", "tooltip")
        .style("background-color", "white")
        .style("border", "solid")
        .style("border-width", "1px")
        .style("border-radius", "5px")
        .style("padding", "10px");
}

// Three function that change the tooltip when user hover / move / leave a cell
export function getMouseOverHandler(tooltip) {
    return () => {
        const singleForecast: Forecast = d3.select(this).datum() as Forecast;

        tooltip
            .html(`Forcasted: ${singleForecast.value.toFixed(2)} vs Actual: ${(singleForecast.actual ?? 0).toFixed(2)}`)
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