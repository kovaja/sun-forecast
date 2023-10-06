export interface ModalResult {
  value: any
}

export type CloseModalCallback = (result: ModalResult) => void
type OpenModalCallback = (component: any, closeCallback: CloseModalCallback) => void

const modal: { openModalCallback: OpenModalCallback } = {
  openModalCallback: undefined
}
export const openModal: OpenModalCallback = (component, closeCallback) => {
  modal.openModalCallback(component, closeCallback)
}

export function registerModal(callback: OpenModalCallback) {
  modal.openModalCallback = callback
}