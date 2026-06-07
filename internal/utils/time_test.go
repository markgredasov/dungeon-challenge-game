package utils

import (
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	tests := []struct {
		name    string
		strTime string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "valid time format",
			strTime: "01:02:03",
			want:    time.Date(0, 1, 1, 1, 2, 3, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "invalid time format",
			strTime: "1:2:3:4",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTime(tt.strTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if result != tt.want {
				t.Errorf("ParseTime() got = %v, expected %v", result, tt.want)
			}
		})
	}
}

func TestTimeToStr(t *testing.T) {
	time := time.Date(0, 1, 1, 1, 2, 3, 0, time.UTC)
	expected := "01:02:03"

	result := TimeToStr(time)
	if result != expected {
		t.Errorf("TimeToStr() got = %v, expected = %v", result, expected)
	}
}

func TestFormatDuration(t *testing.T) {
	startTime := time.Date(0, 1, 1, 1, 0, 0, 0, time.UTC)
	endTime := time.Date(0, 1, 1, 2, 0, 0, 0, time.UTC)

	duration := endTime.Sub(startTime)

	expected := "01:00:00"

	result := FormatDuration(duration)
	if result != expected {
		t.Errorf("FormatDuration() got = %v expected = %v", result, expected)
	}
}
