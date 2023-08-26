<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import { fetchForecast } from '../services/api/forecast.api';
  import { GRAPH_ROOT, plotGraph } from '../services/forecast/graph';

  const halfHourMs = 30 * 60 * 1000
  let windowSize = 8
  let readableWindowSize = ''
  let windowMiddle = new Date()
  let windowLeft = ''
  let windowRight = '';

  async function renderGraph() {
    const response = await fetchForecast(windowSize, windowMiddle)
    if (response) {
      setReadableWindowBoundaries(response.from, response.to)
      plotGraph(response.data)
    }
  }

  let debounceTimeoutId = null

  function debounceRender() {
    if (debounceTimeoutId) {
      clearTimeout(debounceTimeoutId)
    }
    debounceTimeoutId = setTimeout(renderGraph, 400)
  }

  async function updateWindowSizeUp() {
    windowSize += 1
    setReadableWindowSize()
    debounceRender()
  }

  async function updateWindowSizeDown() {
    windowSize -= 1
    setReadableWindowSize()
    debounceRender()
  }

  async function moveWindowMiddleInPast() {
    windowMiddle = new Date(windowMiddle.getTime() - halfHourMs)
    debounceRender()
  }

  async function moveWindowMiddleToFuture() {
    windowMiddle = new Date(windowMiddle.getTime() + halfHourMs)
    debounceRender()
  }

  function setReadableWindowSize() {
    const hrs = windowSize / 2
    readableWindowSize = (Math.floor(hrs) - hrs) === 0 ? `${hrs}.5 hrs` : `${hrs + 0.5}.0 hrs`
  }

  function setReadableWindowBoundaries(from: string, to: string) {
    windowLeft = new Date(from).toLocaleString()
    windowRight = new Date(to).toLocaleString()
  }

  onMount(async () => {
    setReadableWindowSize()
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
            <button on:click={moveWindowMiddleInPast}>&lt;&lt;</button>
            <span class="graph-control_label">
              {windowLeft}
            </span>
        </div>
        <div class="graph-control_variable">
            <button on:click={updateWindowSizeDown}>-</button>
            <span class="graph-control_label">
              Window size: {readableWindowSize}
            </span>
            <button on:click={updateWindowSizeUp}>+</button>
        </div>
        <div class="graph-control_variable">
            <span class="graph-control_label">
              {windowRight}
            </span>
            <button on:click={moveWindowMiddleToFuture}>&gt;&gt;</button>
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
        justify-content: space-evenly;
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