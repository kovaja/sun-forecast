<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import type { ForecastDiff } from '../../types';
  import { DIFF_GRAPH_ROOT, plotGraph } from '../../services/graph/diff/graph';

  export let diffs: ForecastDiff[] = [];

  function renderGraph(newForecasts: ForecastDiff[]) {
    if (diffs.length > 0) {
      plotGraph(newForecasts)
    }
  }

  function reRenderGraph() {
    renderGraph(diffs)
  }

  onMount(async () => {
    window.addEventListener('resize', reRenderGraph)
  })

  onDestroy(() => {
    window.removeEventListener('resize', reRenderGraph)
  })

  $: renderGraph(diffs)
</script>

<div>
  <div class={DIFF_GRAPH_ROOT + ' graph-container'}>Loading...</div>
</div>

<style>
  .graph-container {
    width: 100%;
  }
</style>