<script lang="ts">
  import type { AppEvent } from '../lib/types';
  import { AppEventType } from '../lib/types';
  import { fetchEvents } from '../lib/services/api/event.api';
  import { onMount } from 'svelte';
  import EventsTable from '../lib/components/EventsTable.svelte';
  import EventsControls from '../lib/components/EventsControls.svelte';

  let selectedType: AppEventType = AppEventType.AppError
  let events: AppEvent[] | null = null
  let loading = true

  async function loadEvents() {
    loading = true
    events = await fetchEvents(selectedType)
    loading = false
  }

  async function onSelectedTypeChanged(event: CustomEvent<{ type: AppEventType }>) {
    selectedType = event.detail.type
    await loadEvents()
  }

  onMount(async () => {
    await loadEvents()
  })
</script>

<main>
    <EventsControls on:typeSelected={onSelectedTypeChanged}/>
    <div>{loading ? 'Loading events...' : ''}</div>
    {#if !loading && events}

        <EventsTable events={events}/>
    {/if}
    {#if !loading && !events}
        <div>No events!</div>
    {/if}
</main>
