<script lang="ts">
  import { onMount } from 'svelte';
  import { fetchForecast } from '../lib/services/api/forecast.api';
  import ForecastGraph from '../lib/components/ForecastGraph.svelte';

  onMount(async () => {
    // await fetchForecast()
  })
</script>

<div class="forecast">
    <ForecastGraph/>
</div>