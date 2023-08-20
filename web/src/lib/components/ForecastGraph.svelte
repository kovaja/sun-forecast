<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
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

  onMount(async () => {
    await renderGraph()
    window.addEventListener('resize', renderGraph)
  })

  onDestroy(() => {
    window.removeEventListener('resize', renderGraph)
  })
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