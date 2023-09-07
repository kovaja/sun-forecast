<script lang="ts">
  import ControlsBar from './ControlsBar/ControlsBar.svelte';
  import { ControlsType } from './ControlsBar/types';
  import type { ControlsVariable } from './ControlsBar/types';
  import { createEventDispatcher, onDestroy, onMount } from 'svelte';
  import { isSmallViewport } from '../utils/dom';

  export let windowFrom: string;
  export let windowTo: string;

  const reFetchInterval = 5 * 60 * 1000 // every 5 minutes
  const halfHourMs = 30 * 60 * 1000
  let windowSize = isSmallViewport() ? 6 : 12
  let windowMiddle = new Date()
  let readableWindowSize = ''
  let windowLeft = ''
  let windowRight = ''
  let controls: ControlsVariable[] = []

  function updateControls(wf: string, wt: string) {
    if (!wf || !wt) {
      return
    }
    windowLeft = new Date(wf).toLocaleString()
    windowRight = new Date(wt).toLocaleString()
    const hrs = windowSize / 2
    readableWindowSize = (Math.floor(hrs) - hrs) === 0 ? `${hrs}.5 hrs` : `${hrs + 0.5}.0 hrs`
    recomputeControls()
  }

  const dispatchEvent = createEventDispatcher<{ windowChange: { windowSize: number; windowMiddle: Date } }>()

  function dispatchWindowChange() {
    dispatchEvent('windowChange', {windowSize, windowMiddle})
  }

  function updateWindowSizeUp() {
    windowSize += 1
    dispatchWindowChange()
  }

  function updateWindowSizeDown() {
    windowSize -= 1
    dispatchWindowChange()
  }

  function moveWindowMiddleInPast() {
    windowMiddle = new Date(windowMiddle.getTime() - halfHourMs)
    dispatchWindowChange()
  }

  function moveWindowMiddleToFuture() {
    windowMiddle = new Date(windowMiddle.getTime() + halfHourMs)
    dispatchWindowChange()
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

  let intervalId: number;
  onMount(() => {
    dispatchWindowChange()
    intervalId = window.setInterval(() => {
      windowMiddle = new Date()
      dispatchWindowChange()
    }, reFetchInterval)
  })

  onDestroy(() => {
    window.clearInterval(intervalId)
  })

  $: updateControls(windowFrom, windowTo)
</script>

<div class="forecast-controls">
    {#if controls.length > 0}
        <ControlsBar controls={controls}/>
    {/if}
</div>

<style>
    .forecast-controls {
        margin: 0 0 0 40px;
    }
</style>