<script lang="ts">
  import { AppEventType } from '../types';
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher<{ typeSelected: { type: AppEventType}}>();

  const buttons = [
    {
      type: AppEventType.ForecastConsumed,
      label: 'Forecasts consumed'
    },
    {
      type: AppEventType.ForecastUpdated,
      label: 'Forecasts updated'
    },
    {
      type: AppEventType.AppError,
      label: 'App error'
    }
  ]

  function onTypeSelected(type: AppEventType) {
    console.log('selected type', type)
    dispatch('typeSelected', { type })
  }

  export let selectedType: AppEventType;
</script>

<div class="tabs-controls">
    {#each buttons as button}
        <button on:click={() => onTypeSelected(button.type)}>
            {#if selectedType === button.type}
                <b>{button.label}</b>
            {:else}
                {button.label}
            {/if}
        </button>
    {/each}
</div>


<style>
    .tabs-controls {
        display: flex;
        align-items: center;
        justify-content: flex-start;
    }
</style>
