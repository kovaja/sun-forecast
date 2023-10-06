<script lang="ts">
  import type { ModalResult } from '../../services/modalService';

  export let onClose: (result: ModalResult) => void;

  const options = [
    {label: '6 hours', value: 6},
    {label: '12 hours', value: 12},
    {label: '1 day', value: 24},
    {label: '2 days', value: 48},
    {label: '3 days', value: 72},
  ]

  function onOptionSelected(value: number) {
    onClose({value})
  }
</script>

<div class="window-size">
    <ul>
        {#each options as option}
            <li>
                <button on:click={() => onOptionSelected(option.value)}>
                    {option.label}
                </button>
            </li>
        {/each}
    </ul>
</div>

<style>
    .window-size ul {
        list-style: none;
        margin: 0;
        padding: 0;
        display: flex;
        flex-direction: column;
    }

    .window-size li {
        margin: 0;
        padding: 0;
    }

    .window-size button {
        width: 100%;
        font-size: 26px;
    }
</style>