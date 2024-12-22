<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import { GRAPH_ROOT, plotGraph } from '../../services/graph/forecast/graph';
  import type { Forecast } from '../../types';

  export let forecasts: Forecast[] = [];
  export let graphContainerId: string = GRAPH_ROOT;

  function renderGraph(newForecasts: Forecast[] = forecasts) {
    if (newForecasts.length > 0) {
      plotGraph(newForecasts, graphContainerId)
    }
  }

  function reRenderGraph() {
    renderGraph(forecasts)
  }

  onMount(async () => {
    window.addEventListener('resize', reRenderGraph)
  })

  onDestroy(() => {
    window.removeEventListener('resize', reRenderGraph)
  })

  $: renderGraph(forecasts)
</script>

<div id={graphContainerId} class="graph-container">Loading...</div>

<style>
    .graph-container {
        width: 100%;
        flex: 1 0 auto;
    }
</style>