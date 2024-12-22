<script lang="ts">
  import ForecastGraph from '../lib/components/Forecasts/ForecastGraph.svelte';
  import { onDestroy, onMount } from 'svelte';
  import type { Forecast, ForecastParams } from '../lib/types';
  import { fetchForecast } from '../lib/services/api/forecast.api';
  import { getTodayStartMs, getYesterdayStartMs, mapWindowToForecastParams } from '../lib/utils/date';

  const WEEK_MS = 7 * 24 * 60 * 60 * 1000
  const DAY_HRS_MS = 24 * 60 * 60 * 1000
  const reFetchInterval = 5 * 60 * 1000 // every 5 minutes

  let allForecast: Forecast[] = []
  let todayForecast: Forecast[] = []

  const allForecastId = 'allForecast';
  const todayForecastId = 'todayForecast';

  async function updateForecasts() {
    const yesterdayMs = getYesterdayStartMs();
    const todayMs = getTodayStartMs();
    const allForecastParams: ForecastParams = {
      from: new Date(yesterdayMs).toISOString(),
      to: new Date(yesterdayMs + WEEK_MS).toISOString(),
    }
    const todayForecastParams: ForecastParams = {
      from: new Date(todayMs).toISOString(),
      to: new Date(todayMs + DAY_HRS_MS).toISOString(),
    }

    const [all, today] = await Promise.all([
      fetchForecast(allForecastParams),
      fetchForecast(todayForecastParams)
    ])

    allForecast = all.forecasts;
    todayForecast = today.forecasts;
  }

  let intervalId: number;
  onMount(() => {
    void updateForecasts();

    intervalId = window.setInterval(() => {
      void updateForecasts()
    }, reFetchInterval)
  })

  onDestroy(() => {
    window.clearInterval(intervalId)
  })
</script>

<div class="dashboard">
  <ForecastGraph forecasts={allForecast} graphContainerId={allForecastId}/>
  <hr>
  <ForecastGraph forecasts={todayForecast} graphContainerId={todayForecastId}/>
</div>

<style>
  .dashboard {
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-evenly;
  }
  hr {
    width: 100%;
  }
</style>