<script lang="ts">
  import type { AppEvent } from '../lib/types';
  import { fetchEvents } from '../lib/services/api/event.api';
  import { onMount } from 'svelte';
  import { formatDate } from '../lib/utils/date';

  let events: AppEvent[] = []
  let loading = true

  onMount(async () => {
    events = await fetchEvents()
    loading = false
  })
</script>

<main>
    <h1>App events</h1>
    <div>{loading ? 'Loading events...' : ''}</div>
    <div class="table-container">
        <table>
            <thead>
            <tr>
                <td>Time</td>
                <td>Message</td>
            </tr>
            </thead>
            <tbody>
            {#each events as event}
                <tr>
                    <td>{formatDate(event.timestamp)}</td>
                    <td>{event.message}</td>
                </tr>
            {/each}
            </tbody>
        </table>
    </div>
</main>

<style>
    .table-container {
        height: 80vh;
        overflow: scroll;
    }

    td {
        padding: 10px;
        border: 1px solid;
    }

</style>
