<script lang="ts">
  export let keepLabelVisible: boolean;
  export let label: string;

  let clsx = `controls-bar_label ${keepLabelVisible ? 'controls-bar_label--no-hide' : ''}`
</script>

<span class={clsx}>{label}</span>

<style>
    .controls-bar_label {
        padding: 5px;
        border: 1px solid #395B64;
        font-size: 7px;
        display: none;
    }

    .controls-bar_label--no-hide {
        display: block;
    }

    @media (min-width: 480px) {
        .controls-bar_label {
            display: block;
            font-size: 11px;
        }
    }
</style>