<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import { fetchForecast } from '../services/api/forecast.api';
  import { GRAPH_ROOT, plotGraph } from '../services/forecast/graph';
  import ControlsBar from './ControlsBar/ControlsBar.svelte';
  import type { ControlsVariable } from './ControlsBar/types';
  import { ControlsType } from './ControlsBar/types';

  const halfHourMs = 30 * 60 * 1000
  let windowSize = 8
  let readableWindowSize = ''
  let windowMiddle = new Date()
  let windowLeft = ''
  let windowRight = '';
  let controls: ControlsVariable[] = [];

  function updateControls(from: string, to: string) {
    windowLeft = new Date(from).toLocaleString()
    windowRight = new Date(to).toLocaleString()
    const hrs = windowSize / 2
    readableWindowSize = (Math.floor(hrs) - hrs) === 0 ? `${hrs}.5 hrs` : `${hrs + 0.5}.0 hrs`
    recomputeControls()
  }

  async function renderGraph() {
    const response = await fetchForecast(windowSize, windowMiddle)
    if (response) {
      updateControls(response.from, response.to)
      plotGraph(response.forecasts)
    }
  }

  let debounceTimeoutId = null

  function debounceRender() {
    if (debounceTimeoutId) {
      clearTimeout(debounceTimeoutId)
    }
    debounceTimeoutId = setTimeout(renderGraph, 400)
  }

  function updateWindowSizeUp() {
    windowSize += 1
    debounceRender()
  }

  function updateWindowSizeDown() {
    windowSize -= 1
    debounceRender()
  }

  function moveWindowMiddleInPast() {
    windowMiddle = new Date(windowMiddle.getTime() - halfHourMs)
    debounceRender()
  }

  function moveWindowMiddleToFuture() {
    windowMiddle = new Date(windowMiddle.getTime() + halfHourMs)
    debounceRender()
  }

  function recomputeControls() {
    controls = [
      {
        type: ControlsType.Button,
        sign: '<<',
        label: windowLeft,
        onClick: () => moveWindowMiddleInPast()
      },
      {
        type: ControlsType.Button,
        sign: '-',
        label: 'Window size: ' + readableWindowSize,
        keepLabelVisible: true,
        onClick: () => updateWindowSizeDown()
      },
      {
        type: ControlsType.Button,
        label: '+',
        onClick: () => updateWindowSizeUp()
      },
      {
        type: ControlsType.Button,
        sign: '>>',
        label: windowRight,
        labelPosition: 'left',
        onClick: () => moveWindowMiddleToFuture()
      }
    ]
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
    <div class="graph-container__controls">
        <ControlsBar controls={controls}/>
    </div>
    <div class={GRAPH_ROOT + ' graph-container'}></div>
</div>

<style>
    .graph-container {
        width: 100%;
    }

    .graph-container__controls {
        margin: 0 0 0 40px;
    }
</style>