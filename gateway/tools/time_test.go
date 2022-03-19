package tools

import (
	"testing"
	"time"
)

func TestFormatDate(t *testing.T) {
	now := time.Now()
	if now.Format(timeLayout) != FormatDate(now) {
		t.Error("datetime must be the same ")
	}
}
