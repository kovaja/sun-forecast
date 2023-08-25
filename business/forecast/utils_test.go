package forecast

import (
	"testing"
	"time"
)

func TestGetPeriodEnd(t *testing.T) {
	testCases := []struct {
		name           string
		inputTime      time.Time
		expectedResult time.Time
	}{
		{
			name:           "First half of hour",
			inputTime:      time.Date(2023, time.August, 23, 13, 17, 0, 0, time.UTC),
			expectedResult: time.Date(2023, time.August, 23, 13, 30, 0, 0, time.UTC),
		},
		{
			name:           "Second half of hour",
			inputTime:      time.Date(2023, time.August, 23, 13, 43, 0, 0, time.UTC),
			expectedResult: time.Date(2023, time.August, 23, 14, 00, 0, 0, time.UTC),
		},
		{
			name:           "Just at half of hour",
			inputTime:      time.Date(2023, time.August, 23, 13, 30, 0, 1, time.UTC),
			expectedResult: time.Date(2023, time.August, 23, 14, 00, 0, 0, time.UTC),
		},
		{
			name:           "Just before half of hour",
			inputTime:      time.Date(2023, time.August, 23, 13, 29, 59, 999, time.UTC),
			expectedResult: time.Date(2023, time.August, 23, 13, 30, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			record := HaHistoryRecord{LastChanged: tc.inputTime}
			result := getPeriodEnd(record)

			if !result.Equal(tc.expectedResult) {
				t.Errorf("For: %v, Expected: %v, but got: %v", tc.inputTime, tc.expectedResult, result)
			}
		})
	}
}
