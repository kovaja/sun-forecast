export enum ControlsType {
  Button,
}

interface ControlsBase {
  label: string;
  sign?: string;
  keepLabelVisible?: boolean;
  labelPosition?: 'left' | 'right'
}
interface ControlsButton {
  type: ControlsButton;
  onClick: () => void;
}

export type ControlsVariable = ControlsBase & (ControlsButton)