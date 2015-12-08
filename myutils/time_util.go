package myutils

import (
	"time"
)

type TimeSpanDaily struct {
	// time as int: 19:35:23 will be 193523
	begin int
	until int
}

type TimeSpansDaily []TimeSpanDaily

func (tsd *TimeSpanDaily) WithinSpan(t time.Time) bool {
	timeAsInt := getTimeAsInt(t)
	return timeAsInt >= tsd.begin && timeAsInt <= tsd.until
}

func (tssd *TimeSpansDaily) WithinAnySpan(t time.Time) bool {
	for _, tsd := range *tssd {
		if tsd.WithinSpan(t) {
			return true
		}
	}
	return false
}

func ParseTimeSpans(repeat [][2]int) (tssd TimeSpansDaily) {
	for _, pair := range repeat {
		tssd = append(tssd, TimeSpanDaily{
			begin: pair[0],
			until: pair[1],
		})
	}
	return
}

func getTimeAsInt(t time.Time) int {
	return t.Hour()*10000 + t.Minute()*100 + t.Second()
}
