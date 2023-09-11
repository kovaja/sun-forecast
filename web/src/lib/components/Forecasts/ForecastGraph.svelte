<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import { GRAPH_ROOT, plotGraph } from '../../services/forecast/graph';
  import type { Forecast } from '../../types';

  export let forecasts: Forecast[] = [];

  function renderGraph(newForecasts: Forecast[] = forecasts) {
    if (newForecasts.length > 0) {
      plotGraph(newForecasts)
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

<div>
    <div class={GRAPH_ROOT + ' graph-container'}>Loading...</div>
</div>

<style>
    .graph-container {
        width: 100%;
    }
</style>