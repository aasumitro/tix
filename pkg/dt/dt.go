package dt

import (
	"github.com/aasumitro/tix/common"
	"time"
)

func CurrentDayStartToEnd(currentTime time.Time) (start, end time.Time) {
	startOfDay := time.Date(
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		0, 0, 0, 0, currentTime.Location())
	endOfDay := time.Date(
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		23, 59, 0, 0, currentTime.Location())
	return startOfDay, endOfDay
}

type Weekly struct {
	Start time.Time
	End   time.Time
}

func WeekDayStartToEnd(startDay time.Time) []*Weekly {
	var weekly []*Weekly
	for i := 1; i < common.LastWeekDay; i++ {
		day := startDay.AddDate(0, 0, i-common.LastWeekDay)
		start, end := CurrentDayStartToEnd(day)
		weekly = append(weekly, &Weekly{
			Start: start,
			End:   end,
		})
	}
	start, end := CurrentDayStartToEnd(startDay)
	weekly = append(weekly, &Weekly{
		Start: start,
		End:   end,
	})
	return weekly
}
