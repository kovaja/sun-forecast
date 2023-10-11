<script lang="ts">
  import ControlsBar from '../ControlsBar/ControlsBar.svelte';
  import type { ControlsVariable } from '../ControlsBar/types';
  import { ControlsType } from '../ControlsBar/types';
  import { createEventDispatcher, onDestroy, onMount } from 'svelte';

  let windowSize = 5
  let controls: ControlsVariable[] = []

  const dispatchEvent = createEventDispatcher<{ windowChange: { windowSize: number } }>()

  function dispatchWindowChange() {
    dispatchEvent('windowChange', { windowSize })
    recomputeControls()
  }

  function updateWindowSizeUp() {
    windowSize += 1
    dispatchWindowChange()
  }

  function updateWindowSizeDown() {
    windowSize -= 1
    dispatchWindowChange()
  }

  function recomputeControls() {
    controls = [
      {
        type: ControlsType.Group,
        label: `Window size: ${windowSize + 1} day(s)`,
        centerFieldType: 'text',
        keepLabelVisible: true,
        leftButton:{
          sign: '-',
          onClick: () => updateWindowSizeDown()
        },
        rightButton: {
          sign: '+',
          onClick: () => updateWindowSizeUp()
        }
      }
    ]
  }

  onMount(() => {
    dispatchWindowChange()
  })
</script>

<div class="diff-controls">
    {#if controls.length > 0}
        <ControlsBar controls={controls} align="center"/>
    {/if}
</div>
