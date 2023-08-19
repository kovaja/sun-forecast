<script lang="ts">
    import * as d3 from "d3";
    import { onMount } from 'svelte';
    import type { Forecast } from '../types';
    import { fetchForecast } from '../services/api/forecast.api';

    // const mockData = [
    //     {
    //         "id": 251,
    //         "periodEnd": "2023-08-19T14:00:00Z",
    //         "value": 2431.4,
    //         "actual": 2691.25
    //     },
    //     {
    //         "id": 252,
    //         "periodEnd": "2023-08-19T14:30:00Z",
    //         "value": 1932,
    //         "actual": 1963.7142857142858
    //     },
    //     {
    //         "id": 253,
    //         "periodEnd": "2023-08-19T15:00:00Z",
    //         "value": 1550.4,
    //         "actual": 1887.3164556962026
    //     },
    //     {
    //         "id": 254,
    //         "periodEnd": "2023-08-19T15:30:00Z",
    //         "value": 1070.3,
    //         "actual": null
    //     }
    // ]

    function getContainerWidth(): number {
        return document.querySelector('.graph-root').getBoundingClientRect().width
    }


    const min30Ms = 30 * 60 * 1000
    const colPad = 15;

    function getPeriodStart(periodEnd: string): Date {
        const endMs = new Date(periodEnd).getTime()
        return new Date(endMs - min30Ms)
    }
    function getPeriodMiddle(periodEnd: string): Date {
        const endMs = new Date(periodEnd).getTime()
        return new Date(endMs - min30Ms/2)
    }

    function plotGraph(data: Forecast[]) {


// set the dimensions and margins of the graph

        const margin = {top: 10, left: 35, right: 5, bottom: 70};
        const width = getContainerWidth();
        const height = 600 - margin.top - margin.bottom;
        const rightEdge = width - margin.left - margin.right
        const bottomEdge = height - 50

        const svg = d3.select('.graph-root')
            .append("svg")
            .attr("width", width)
            .attr("height", height)
            .append("g")
            .attr(
                "transform",
                "translate(" + margin.left + "," + margin.top + ")"
            );

        // X axis

        const x = d3.scaleTime()
            .range([0, rightEdge])
            .domain([
                getPeriodStart(data[0].periodEnd),
                new Date(data[data.length - 1].periodEnd)
            ])


        // append to graph
        svg.append("g")
            .attr("transform", "translate(0," + bottomEdge + ")")
            .call(d3.axisBottom(x))
            .selectAll("text")
            .attr("transform", "translate(-10,0)rotate(-45)")
            .style("text-anchor", "end");


        // Add Y axis
        const y = d3.scaleLinear()
            .domain([0, d3.max(data, d => d3.max([d.value, d.actual])) * 1.20])
            .range([bottomEdge, 0]);
        // append to graph
        svg.append("g")
            .call(d3.axisLeft(y));




        // ----------------
        // Create a tooltip
        // ----------------
        const tooltip = d3.select('.graph-root')
            .append("div")
            .style("opacity", 0)
            .attr("class", "tooltip")
            .style("background-color", "white")
            .style("border", "solid")
            .style("border-width", "1px")
            .style("border-radius", "5px")
            .style("padding", "10px")

        // Three function that change the tooltip when user hover / move / leave a cell
        const mouseover = function() {
            const singleForecast: Forecast = d3.select(this).datum();

            tooltip
                .html(`Forcasted: ${singleForecast.value.toFixed(2)} vs Actual: ${(singleForecast.actual ?? 0).toFixed(2)}`)
                .style("opacity", 1)
        }
        // const mousemove = function(d) {
        //     tooltip
        //         .style("left", (d3.mouse(this)[0]+90) + "px") // It is important to put the +90: other wise the tooltip is exactly where the point is an it creates a weird effect
        //         .style("top", (d3.mouse(this)[1]) + "px")
        // }
        const mouseleave = function(d) {
            tooltip
                .style("opacity", 0)
        }

        svg.selectAll("mybar")
            .data(data)
            .enter()
            .append("rect")
            .attr("x", function (d) {
                return colPad / 2 + x(getPeriodStart(d.periodEnd));
            })
            .attr("y", function (d) {
                return y(d.value);
            })
            .attr("width", (rightEdge / data.length) - colPad)
            .attr("height", function (d) {
                return bottomEdge - y(d.value);
            })
            .attr("fill", "#69b3a2")

            .on("mouseover", mouseover)
            // .on("mousemove", mousemove)
            .on("mouseleave", mouseleave)

        svg.selectAll("mybar2")
            .data(data)
            .enter()
            .append("rect")
            .attr("x", function (d) {
                return colPad / 2 + x(getPeriodStart(d.periodEnd));
            })
            .attr("y", function (d) {
                return y(d.actual ?? 0);
            })
            .attr("width", (rightEdge / data.length) - colPad)
            .attr("height", function (d) {
                return bottomEdge - y(d.actual ?? 0);
            })
            .attr("fill", "rgba(194,137,33,0.29)")
            .attr("stroke", "#c28921")
            .on("mouseover", mouseover)
            // .on("mousemove", mousemove)
            .on("mouseleave", mouseleave)


        svg.selectAll("successRate")
            .data(data)
            .enter()
            .append("text")
            .attr("x", function (d) {
                return colPad / 2 + x(getPeriodStart(d.periodEnd));
            })
            .attr("y", function (d) {
                return d3.min([y(d.actual ?? 0), y(d.value)]) - 10;
            })
            .text(d => {
                if (d.value === 0 || d.actual === null) {
                    return ""
                }
                return (-1*(100-(d.actual/d.value)*100)).toFixed(2) + '%'
            })
            .attr("fill", "red")
            .attr("font-size", "12px")
    }

    onMount(async () => {
        const data = await fetchForecast()
        if (data) {
            plotGraph(data)
        }
    })
</script>

<div>
    <div class="graph-root"></div>
</div>

<style>
    .graph-root {
        width: 100%;
    }
</style>