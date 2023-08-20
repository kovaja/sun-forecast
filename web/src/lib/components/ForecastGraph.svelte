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
        <div class="graph-control_variable">
            <button on:click={updateWindowSizeDown}>-</button>
            <span class="graph-control_label">
              Window size: {windowSize}
            </span>
            <button on:click={updateWindowSizeUp}>+</button>
        </div>
    </div>
    <div class={GRAPH_ROOT + ' graph-container'}></div>
</div>

<style>
    .graph-container {
        width: 100%;
    }
    .graph-control {
        margin: 0 0 0 40px;
        padding: 2px 0;
        background-color: #A5C9CA;
        color: #395B64;
        display: flex;
        justify-content: center;
    }
    .graph-control_variable {
        display: flex;
        align-items: center;
    }
    .graph-control_label {
        padding: 5px;
        border: 1px solid #395B64;
        font-size: 11px;
    }
</style>