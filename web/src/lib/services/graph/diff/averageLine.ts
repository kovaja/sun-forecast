import type { ForecastDiff } from '../../../types';

export function computeAverageSeries(diffs: ForecastDiff[]): number[] {
  const diffsCount = diffs[0].diffs.length
  const avgDiff = []
  for (let i = 0; i < diffsCount; i++) {
    let sum = 0
    for (let j = 0; j < diffs.length; j++) {
      sum += diffs[j].diffs[i]
    }
    avgDiff.push(sum / diffs.length)
  }
  return avgDiff
}
