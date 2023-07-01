package dt_test

import (
	"github.com/aasumitro/tix/pkg/dt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_CurrentDayStartToEnd(t *testing.T) {
	testCases := []struct {
		input    time.Time
		expected struct {
			start time.Time
			end   time.Time
		}
	}{
		{
			input: time.Date(2023, 6, 3, 12, 34, 56, 0, time.UTC),
			expected: struct {
				start time.Time
				end   time.Time
			}{
				start: time.Date(2023, 6, 3, 0, 0, 0, 0, time.UTC),
				end:   time.Date(2023, 6, 3, 23, 59, 0, 0, time.UTC),
			},
		},
	}

	for _, testCase := range testCases {
		start, end := dt.CurrentDayStartToEnd(testCase.input)
		assert.Equal(t, start, testCase.expected.start)
		assert.Equal(t, end, testCase.expected.end)
	}
}

func Test_WeekDayStartToEnd(t *testing.T) {
	startDay := time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)
	expectedWeekly := []*dt.Weekly{
		{
			Start: time.Date(2023, 5, 26, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2023, 5, 26, 23, 59, 0, 0, time.UTC),
		},
		{
			Start: time.Date(2023, 5, 27, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2023, 5, 27, 23, 59, 0, 0, time.UTC),
		},
		{
			Start: time.Date(2023, 5, 28, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2023, 5, 28, 23, 59, 0, 0, time.UTC),
		},
		{
			Start: time.Date(2023, 5, 29, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2023, 5, 29, 23, 59, 0, 0, time.UTC),
		},
		{
			Start: time.Date(2023, 5, 30, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2023, 5, 30, 23, 59, 0, 0, time.UTC),
		},
		{
			Start: time.Date(2023, 5, 31, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2023, 5, 31, 23, 59, 0, 0, time.UTC),
		},
		{
			Start: time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2023, 6, 1, 23, 59, 0, 0, time.UTC),
		},
	}
	weekly := dt.WeekDayStartToEnd(startDay)
	assert.Equal(t, len(weekly), len(expectedWeekly))
	for i, w := range weekly {
		expected := expectedWeekly[i]
		assert.Equal(t, w.Start, expected.Start)
		assert.Equal(t, w.End, expected.End)
	}
}
