<script lang="ts">
  import type { ControlsVariable } from './types';
  import { ControlsType } from './types';
  import { beforeUpdate } from 'svelte';

  export let controls: ControlsVariable[];

  beforeUpdate(() => {
    console.log('update controls', controls)
  })
</script>

<div class="controls-bar">
    {#each controls as control}
        <div class="controls-bar_variable">
            {#if control.type === ControlsType.Button}
                {#if control.sign && control.labelPosition === 'left'}
                    {@const clsx = `controls-bar_label ${control.keepLabelVisible ? 'controls-bar_label--no-hide' : ''}`}
                    <span class={clsx}>
                        {control.label}
                    </span>
                {/if}
                <button on:click={control.onClick}>
                    {control.sign ?? control.label}
                </button>
                {#if control.sign && (!control.labelPosition || control.labelPosition === 'right')}
                    {@const clsx = `controls-bar_label ${control.keepLabelVisible ? 'controls-bar_label--no-hide' : ''}`}
                    <span class={clsx}>
                        {control.label}
                    </span>
                {/if}

            {/if}
        </div>
    {/each}
</div>
<style>
    .controls-bar {
        padding: 2px 0;
        background-color: #A5C9CA;
        color: #395B64;
        display: flex;
        justify-content: space-between;
    }

    .controls-bar_variable {
        display: flex;
        align-items: center;
    }

    .controls-bar_label {
        padding: 5px;
        border: 1px solid #395B64;
        font-size: 11px;
        display: none;
    }

    .controls-bar_label--no-hide {
        display: block;
    }


    @media (min-width: 480px) {
        .controls-bar_label {
            display: block;
        }
    }
</style>