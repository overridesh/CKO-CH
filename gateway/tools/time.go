package tools

import (
	"time"
)

const timeLayout string = "2006-01-02T15:04:05.000Z"

func FormatDate(t time.Time) string {
	return t.Format(timeLayout)
}
