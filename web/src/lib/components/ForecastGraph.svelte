<script lang="ts">
  import { onMount } from 'svelte';
  import { fetchForecast } from '../services/api/forecast.api';
  import { GRAPH_ROOT, plotGraph } from '../services/forecast/graph';

  let windowSize = 12
  let windowMiddle = new Date()

  async function renderGraph() {
    const data = await fetchForecast(windowSize, windowMiddle)
    if (data) {
      plotGraph(data)
    }
  }

  async function updateWindowSizeUp() {
    windowSize += 1
    await renderGraph()
  }

  async function updateWindowSizeDown() {
    windowSize -= 1
    await renderGraph()
  }

  onMount(renderGraph)
</script>

<div>
    <div class="graph-control">
        <div class="graph-control_window-size">
            <button on:click={updateWindowSizeDown}>&lt;</button>
            Window size: {windowSize}
            <button on:click={updateWindowSizeUp}>&gt;</button>
        </div>
    </div>
    <div class={GRAPH_ROOT + ' graph-container'}></div>
</div>

<style>
    .graph-container {
        width: 100%;
    }
</style>