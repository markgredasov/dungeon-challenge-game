package utils

import (
	"fmt"
	"time"
)

func ParseTime(strTime string) (time.Time, error) {
	t, err := time.Parse("15:04:05", strTime)
	if err != nil {
		return time.Time{}, fmt.Errorf("parsing time: %w", err)
	}
	return t, nil
}

func TimeToStr(t time.Time) string {
	return t.Format("15:04:05")
}

func FormatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
