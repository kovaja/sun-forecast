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

export function computeMedianSeries(diffs: ForecastDiff[]): number[] {
  const diffsCount = diffs[0].diffs.length;
  const medianDiff = [];

  for (let i = 0; i < diffsCount; i++) {
    const valuesAtPoint = [];

    for (let j = 0; j < diffs.length; j++) {
      valuesAtPoint.push(diffs[j].diffs[i]);
    }

    // Sort the values in ascending order
    valuesAtPoint.sort((a, b) => a - b);

    if (valuesAtPoint.length % 2 === 0) {
      // If there is an even number of values, take the average of the middle two
      const middleIndex = valuesAtPoint.length / 2;
      const median = (valuesAtPoint[middleIndex - 1] + valuesAtPoint[middleIndex]) / 2;
      medianDiff.push(median);
    } else {
      // If there is an odd number of values, take the middle one
      const middleIndex = Math.floor(valuesAtPoint.length / 2);
      medianDiff.push(valuesAtPoint[middleIndex]);
    }
  }

  return medianDiff;
}