package scrappy

import "time"

const THRESHOLD = time.Minute

func IsValidPeriod(unix int) bool {
	now := Now()

	return unix == now
}

func Now() int {
	period := time.Now()
	period = time.Date(period.Year(), period.Month(), period.Day(), period.Hour(), period.Minute(), 0, 0, time.UTC)
	unixPeriod := period.Unix()

	return int(unixPeriod)
}
