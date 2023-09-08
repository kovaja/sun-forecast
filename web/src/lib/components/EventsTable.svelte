<script lang="ts">
  import type { AppEvent } from "../types";
  import { formatDate } from '../utils/date';

  export let events: AppEvent[] = []
</script>

<div class="table-container">
    <ul>
        {#each events as event}
            <li>
                <div class="table-container__date">{formatDate(event.timestamp, true)}</div>
                <div class="table-container__message">{event.message}</div>
            </li>
        {/each}
    </ul>
</div>

<style>
    .table-container {
        height: 85vh;
        overflow: auto;
        font-size: 11px;
    }

    ul {
        list-style: none;
        padding: 0;
        margin: 0;
    }

    li {
        padding: 10px;
        border: 1px solid;
        display: flex;
    }

    .table-container__date {
        flex: 0 1 auto;
        padding-right: 10px;
    }

    .table-container__message {
        flex: 1 1 auto;
        padding-left: 10px;
        border-left: 1px solid;
    }

</style>