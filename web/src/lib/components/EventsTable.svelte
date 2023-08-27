<script lang="ts">
  import type { AppEvent } from "../types";
  import { formatDate } from '../utils/date';

  export let events: AppEvent[] = []
</script>

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
                <td>{formatDate(event.timestamp, true)}</td>
                <td>{event.message}</td>
            </tr>
        {/each}
        </tbody>
    </table>
</div>

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