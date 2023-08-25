package forecast

import (
	"reflect"
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

func TestAppendUpdate(t *testing.T) {
	newUpdateMock := ForecastUpdate{
		PeriodEnd:    time.Date(2023, time.August, 25, 21, 30, 00, 00, time.UTC),
		Actual:       1354,
		ActualCount:  120,
		LastActualAt: time.Date(2023, time.August, 25, 21, 28, 53, 00, time.UTC),
	}

	testCases := []struct {
		name               string
		existingUpdates    []ForecastUpdate
		updated            bool
		expectedUpdatesLen int
		expectedUpdates    []ForecastUpdate
	}{
		{
			name:               "Will not modify the given collection when updated is false",
			existingUpdates:    []ForecastUpdate{},
			updated:            false,
			expectedUpdatesLen: 0,
			expectedUpdates:    []ForecastUpdate{},
		},
		{
			name:               "Will add given update to updates when updated is true",
			existingUpdates:    []ForecastUpdate{},
			updated:            true,
			expectedUpdatesLen: 1,
			expectedUpdates:    []ForecastUpdate{newUpdateMock},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := appendUpdate(tc.existingUpdates, tc.updated, newUpdateMock)

			if len(result) != tc.expectedUpdatesLen {
				t.Errorf("%s: Expected len: %d, but got: %d", tc.name, tc.expectedUpdatesLen, len(result))
			}
			if !reflect.DeepEqual(result, tc.expectedUpdates) {
				t.Errorf("%s: Result is not equal with expected collection", tc.name)
			}
		})
	}
}

func TestComputeUpdatesForOneRecordAndAlreadyTouchedForecast(t *testing.T) {
	periodEnd := time.Date(2023, time.August, 25, 23, 00, 00, 00, time.UTC)

	haRecordMock := HaHistoryRecord{
		LastChanged: time.Date(2023, time.August, 25, 22, 38, 01, 00, time.UTC),
		State:       "163",
	}

	loaderMock := func(t time.Time) *Forecast {
		actual := float64(2155)
		lastActualAt := time.Date(2023, time.August, 25, 22, 36, 00, 00, time.UTC)
		return &Forecast{
			Id:           1,
			PeriodEnd:    periodEnd,
			Value:        2645,
			Actual:       &actual,
			ActualCount:  5,
			LastActualAt: &lastActualAt,
		}
	}

	updates := ComputeUpdates(loaderMock, []HaHistoryRecord{haRecordMock})

	if len(updates) != 1 {
		t.Errorf("Expected 1 updates, got %d", len(updates))
		return
	}

	expectedUpdate := ForecastUpdate{
		PeriodEnd:    periodEnd,
		Actual:       (2155*5 + 163) / 6, //cummulative average
		ActualCount:  5 + 1,
		LastActualAt: haRecordMock.LastChanged,
	}

	if !reflect.DeepEqual(updates[0], expectedUpdate) {
		t.Errorf("Unexpected data in update. Expected %v, got %v", expectedUpdate, updates[0])
	}
}

func TestComputeUpdatesForTwoRecordsAndUntouchedForecast(t *testing.T) {
	periodEnd := time.Date(2023, time.August, 25, 23, 00, 00, 00, time.UTC)

	haRecordsMock := []HaHistoryRecord{
		{
			LastChanged: time.Date(2023, time.August, 25, 22, 38, 01, 00, time.UTC),
			State:       "163",
		},
		{
			LastChanged: time.Date(2023, time.August, 25, 22, 39, 01, 00, time.UTC),
			State:       "250",
		},
	}

	loaderMock := func(t time.Time) *Forecast {
		return &Forecast{
			Id:           1,
			PeriodEnd:    periodEnd,
			Value:        2645,
			Actual:       nil,
			ActualCount:  0,
			LastActualAt: nil,
		}
	}

	updates := ComputeUpdates(loaderMock, haRecordsMock)

	if len(updates) != 1 {
		t.Errorf("Expected 1 updates, got %d", len(updates))
		return
	}

	expectedUpdate := ForecastUpdate{
		PeriodEnd:    periodEnd,
		Actual:       206.5,
		ActualCount:  2,
		LastActualAt: haRecordsMock[1].LastChanged,
	}

	if !reflect.DeepEqual(updates[0], expectedUpdate) {
		t.Errorf("Unexpected data in update. Expected %v, got %v", expectedUpdate, updates[0])
	}
}
