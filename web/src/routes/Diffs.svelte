<script lang="ts">
  import ForecastGraph from '../lib/components/ForecastGraph.svelte';
  import ForecastControls from '../lib/components/ForecastControls.svelte';
  import { fetchDiffs} from '../lib/services/api/forecast.api';
  import type { ForecastDiff } from '../lib/types';
  import DiffsGraph from '../lib/components/Diffs/DiffsGraph.svelte';
  import { onMount } from 'svelte';

  // triggered and initialized by controls component
  let windowSize: number;
  let windowMiddle: Date;

  let windowFrom = ''
  let windowTo = ''

  let diffs: ForecastDiff[] = [];

  async function _fetchDiffs() {
    diffs = await fetchDiffs()
  }

  let debounceTimeoutId = null

  function debounceFetch() {
    if (debounceTimeoutId) {
      clearTimeout(debounceTimeoutId)
    }
    debounceTimeoutId = setTimeout(_fetchDiffs, 400)
  }

  function onWindowChange({detail: {windowSize: ws, windowMiddle: wm}}: CustomEvent) {
    windowSize = ws;
    windowMiddle = wm;
    debounceFetch()
  }

  onMount(debounceFetch)
</script>

<div class="diffs">
  <DiffsGraph diffs={diffs} />
</div>