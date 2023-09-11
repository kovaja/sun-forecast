<script lang="ts">
  import { fetchDiffs} from '../lib/services/api/forecast.api';
  import type { ForecastDiff } from '../lib/types';
  import DiffsGraph from '../lib/components/Diffs/DiffsGraph.svelte';
  import DiffsControls from '../lib/components/Diffs/DiffsControls.svelte';

  // triggered and initialized by controls component
  let windowSize: number;

  let diffs: ForecastDiff[] = [];

  async function _fetchDiffs() {
    diffs = await fetchDiffs(windowSize)
  }

  let debounceTimeoutId = null

  function debounceFetch() {
    if (debounceTimeoutId) {
      clearTimeout(debounceTimeoutId)
    }
    debounceTimeoutId = setTimeout(_fetchDiffs, 400)
  }

  function onWindowChange({detail: {windowSize: ws}}: CustomEvent) {
    windowSize = ws;
    debounceFetch()
  }
</script>

<div class="diffs">
  <DiffsControls on:windowChange={onWindowChange} />
  <DiffsGraph diffs={diffs} />
</div>