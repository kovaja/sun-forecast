<script lang="ts">
  import { onMount } from 'svelte';
  import { registerModal } from '../../services/modalService';
  import type { ModalResult } from '../../services/modalService'
  import type { Modal } from './types';

  let opened: boolean = false;
  let component: Modal;
  let closeCallback: (result: ModalResult) => void

  function openModal(_component: any, _closeCallback: (result: ModalResult) => void) {
    console.log('component open modal', _closeCallback)
    closeCallback = _closeCallback
    component = _component
    opened = true
  }

  function tearDownModal() {
    opened = false
    component = undefined
    closeCallback = undefined
  }

  function onCloseModal(result: ModalResult) {
    closeCallback(result)
    tearDownModal()
  }

  function onBlur() {
    tearDownModal()
  }

  onMount(() => {
    registerModal(openModal)
  })



</script>

{#if opened}
  <dialog on:blur={onBlur}>
    <svelte:component this={component} onClose={onCloseModal} />
  </dialog>
{/if}

<style>
  dialog {
    display: block;
    position: fixed;
    top: 50%;
    transform: translateY(-50%);
    padding: 0;
  }
</style>
