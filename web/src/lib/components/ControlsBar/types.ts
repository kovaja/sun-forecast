export enum ControlsType {
  Button,
  Group
}

interface ControlsBase {
  label: string;
  sign?: string;
  keepLabelVisible?: boolean;

}
interface ControlsButton {
  labelPosition?: 'left' | 'right'
  type: ControlsType.Button;
  onClick: () => void;
}
interface ControlsGroup {
  type: ControlsType.Group,
  leftButton: {
    sign:string;
    onClick: () => void;
  };
  rightButton: {
    sign:string;
    onClick: () => void;
  }
}


export type ControlsVariable = ControlsBase & (ControlsButton | ControlsGroup)