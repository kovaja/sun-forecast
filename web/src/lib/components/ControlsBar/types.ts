export type ControlsAlign = 'evenly' | 'center'

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

type ControlsGroup = {
  type: ControlsType.Group;
  leftButton: {
    sign: string;
    onClick: () => void;
  };
  rightButton: {
    sign: string;
    onClick: () => void;
  };
} & ({
  centerFieldType: 'button',
  onCenterFieldClick: () => void;
} | {
  centerFieldType: 'text',
  onCenterFieldClick: never;
});


export type ControlsVariable = ControlsBase & (ControlsButton | ControlsGroup)