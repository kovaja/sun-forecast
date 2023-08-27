<script lang="ts">
  import ForecastGraph from '../lib/components/ForecastGraph.svelte';
  import ForecastControls from '../lib/components/ForecastControls.svelte';
  import { fetchForecast } from '../lib/services/api/forecast.api';
  import type { Forecast } from '../lib/types';

  // triggered and initialized by controls component
  let windowSize: number;
  let windowMiddle: Date;

  let windowFrom = ''
  let windowTo = ''

  let forecasts: Forecast[] = [];

  async function fetchForecasts() {
    const response = await fetchForecast(windowSize, windowMiddle)
    if (response) {
      windowFrom = response.from
      windowTo = response.to
      forecasts = response.forecasts
    }
  }

  let debounceTimeoutId = null

  function debounceFetch() {
    if (debounceTimeoutId) {
      clearTimeout(debounceTimeoutId)
    }
    debounceTimeoutId = setTimeout(fetchForecasts, 400)
  }

  function onWindowChange({detail: {windowSize: ws, windowMiddle: wm}}: CustomEvent) {
    windowSize = ws;
    windowMiddle = wm;
    debounceFetch()
  }
</script>

<div class="forecast">
    <ForecastControls on:windowChange={onWindowChange} windowFrom={windowFrom} windowTo={windowTo}/>
    <ForecastGraph forecasts={forecasts}/>
</div>