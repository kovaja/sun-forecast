<script lang="ts">
    import { onMount } from 'svelte';
    import { fetchForecast } from '../services/api/forecast.api';
    import { GRAPH_ROOT, plotGraph } from '../services/forecast/graph';

    onMount(async () => {
        const data = await fetchForecast()
        if (data) {
            plotGraph(data)
        }
    })
</script>

<div>
    <div class={GRAPH_ROOT + ' graph-container'}></div>
</div>

<style>
    .graph-container {
        width: 100%;
    }
</style>