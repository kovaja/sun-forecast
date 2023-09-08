import type { BaseType } from 'd3';
import type { D3Selection } from '../../types.d3';

export function addAttributes<T extends BaseType>(
  selection: D3Selection<T>,
  attrs: Record<string, Parameters<D3Selection<T>['attr']>[1]>
): D3Selection<T> {
  Object.entries(attrs).forEach(([attr, value]) => {
    selection.attr(attr, value)
  })

  return selection
}