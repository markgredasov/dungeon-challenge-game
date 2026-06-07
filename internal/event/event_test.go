package event

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestParseEvents(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantLen int
		wantErr bool
	}{
		{
			name:    "valid events content",
			content: "[14:00:00] 1 1\n[14:00:00] 2 1",
			wantLen: 2,
			wantErr: false,
		},
		{
			name:    "invalid events content",
			content: "invalid event content\n[14:00:00] 2 1",
			wantLen: 1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "events-*.txt")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.WriteString(tt.content); err != nil {
				t.Fatal(err)
			}
			tmpFile.Close()

			filename := tmpFile.Name()
			events, err := ParseEvents(&filename)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseEvents() error = %v wantErr = %v", err, tt.wantErr)
				return
			}

			if len(*events) != tt.wantLen {
				t.Fatalf("ParseEvents() got (length) %v wanted (lenght) %v", len(*events), tt.wantLen)
			}
		})
	}
}

func TestParseEvent(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		want    Event
		wantErr bool
	}{
		{
			name: "valid event",
			line: "[14:00:00] 1 1",
			want: Event{
				Time:        time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC),
				PlayerID:    1,
				ExtraParams: "",
				Type:        EventRegister,
			},
			wantErr: false,
		},
		{
			name:    "not enough params",
			line:    "[14:00:00] 1",
			want:    Event{},
			wantErr: true,
		},
		{
			name:    "invalid event time",
			line:    "[14:00] 1 1",
			want:    Event{},
			wantErr: true,
		},
		{
			name:    "invalid playerID",
			line:    "[14:00:00] a",
			want:    Event{},
			wantErr: true,
		},
		{
			name:    "invalid eventType",
			line:    "[14:00:00] 1 a",
			want:    Event{},
			wantErr: true,
		},
		{
			name: "damage event",
			line: "[14:00:00] 1 11 50",
			want: Event{
				Time:        time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC),
				PlayerID:    1,
				ExtraParams: "50",
				Type:        EventDamage,
			},
			wantErr: false,
		},
		{
			name: "cannot continue event",
			line: "[14:00:00] 1 9 too tired (?)",
			want: Event{
				Time:        time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC),
				PlayerID:    1,
				ExtraParams: "too tired (?)",
				Type:        EventCannotContinue,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseEvent(tt.line)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseEvent() error = %v wantErr = %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(result, tt.want) {
				t.Fatalf("parseEvent() got = %v expected %v", result, tt.want)
			}
		})
	}
}
