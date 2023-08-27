<script lang="ts">
  import { AppEventType } from '../types';
  import { createEventDispatcher } from 'svelte';
  import ControlsBar from './ControlsBar/ControlsBar.svelte';
  import type { ControlsVariable } from './ControlsBar/types';
  import { ControlsType } from './ControlsBar/types';

  const dispatch = createEventDispatcher<{ typeSelected: { type: AppEventType}}>();
  function onTypeSelected(type: AppEventType) {
    dispatch('typeSelected', { type })
  }

  const controls: ControlsVariable[] = [
    {
      type: ControlsType.Button,
      keepLabelVisible: true,
      label: 'Forecasts consumed',
      onClick: () => onTypeSelected(AppEventType.ForecastConsumed)
    },
    {
      type: ControlsType.Button,
      keepLabelVisible: true,
      label: 'Forecasts updated',
      onClick: () => onTypeSelected(AppEventType.ForecastUpdated)
    },
    {
      type: ControlsType.Button,
      keepLabelVisible: true,
      label: 'App error',
      onClick: () => onTypeSelected(AppEventType.AppError)
    }
  ]
</script>

<div class="tabs-controls">
    <ControlsBar controls={controls}/>
</div>


<style>
    .tabs-controls {
        display: flex;
        align-items: center;
        justify-content: flex-start;
    }
</style>
