import type { BaseType, ScaleLinear, ScaleTime, Selection } from 'd3';

export type D3Selection<T extends BaseType> = Selection<T, any, any, any>
export type D3TimeScale = ScaleTime<number, number>
export type D3LinearScale = ScaleLinear<number, number>