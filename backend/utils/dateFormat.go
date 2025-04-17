package utils

import (
	"strconv"
	"time"
)

// DateFromat converts a given time into a human-readable format representing
// the time difference in seconds, minutes, hours, days, or weeks.
func DateFromat(date time.Time) string {
	diff := time.Since(date)
	switch {
	case diff.Seconds() < 60:
		return strconv.Itoa(int(diff.Seconds())) + "s"
	case diff.Minutes() < 60:
		return strconv.Itoa(int(diff.Minutes())) + "min"
	case diff.Hours() < 24:
		return strconv.Itoa(int(diff.Hours())) + "h"
	case diff.Hours() < 24*7:
		return strconv.Itoa(int(diff.Hours()/24)) + "j"
	default:
		return strconv.Itoa(int(diff.Hours()/(24*7))) + "week"
	}
}
