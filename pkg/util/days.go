package util

import (
	"github.com/cryptorities/cryptopass/pkg/app"
	"time"
)

/**
	Alex Shvid
 */

func DaysOffset(t *time.Time) int {
	return int( t.Sub(app.StartDate) / app.TimeDay )
}

func ParseDaysOffset(days int) time.Time {
	t := app.StartDate
	return t.Add(app.TimeDay * time.Duration(days))
}
