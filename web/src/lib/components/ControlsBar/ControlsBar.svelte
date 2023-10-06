<script lang="ts">
  import type { ControlsAlign, ControlsVariable } from './types';
  import { ControlsType } from './types';
  import ControlsLabel from './ControlsLabel.svelte';

  export let controls: ControlsVariable[];
  export let align: ControlsAlign = 'evenly'
</script>

<div class={'controls-bar controls-bar--' + align}>
    {#each controls as control}
        <div class="controls-bar_variable">
            {#if control.type === ControlsType.Button}
                {#if control.sign && control.labelPosition === 'left'}
                    <ControlsLabel keepLabelVisible={control.keepLabelVisible} label={control.label}/>
                {/if}
                <button on:click={control.onClick}>
                    {control.sign ?? control.label}
                </button>
                {#if control.sign && (!control.labelPosition || control.labelPosition === 'right')}
                    <ControlsLabel keepLabelVisible={control.keepLabelVisible} label={control.label}/>
                {/if}

            {:else if control.type === ControlsType.Group}
                <button on:click={control.leftButton.onClick}>
                    {control.leftButton.sign}
                </button>
                {#if control.centerFieldType === 'text'}
                  <ControlsLabel keepLabelVisible={control.keepLabelVisible} label={control.label}/>
                {:else if control.centerFieldType === 'button'}
                  <button on:click={control.onCenterFieldClick}>
                    {control.label}
                  </button>
                {/if}
                <button on:click={control.rightButton.onClick}>
                    {control.rightButton.sign}
                </button>
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

        width: 100%;
    }

    .controls-bar--evenly {
        justify-content: space-between;
    }
    .controls-bar--center {
        justify-content: center;
    }

    .controls-bar_variable {
        display: flex;
        align-items: center;
    }
</style>