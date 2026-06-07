package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		jsonContent string
		want        *Config
		wantErr     bool
	}{
		{
			name: "valid config",
			jsonContent: `{
				"Floors": 5,
				"Monsters": 3,
				"OpenAt": "14:05:00",
				"Duration": 2
			}`,
			want: &Config{
				Floors:   5,
				Monsters: 3,
				OpenAt:   "14:05:00",
				Duration: 2,
			},
			wantErr: false,
		},
		{
			name: "invalid JSON",
			jsonContent: `{
				"Floors": "invalid",
			}`,
			wantErr: true,
		},
		{
			name:        "missing fields",
			jsonContent: `{}`,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "config-temp.json")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.WriteString(tt.jsonContent); err != nil {
				t.Fatal(err)
			}
			tmpFile.Close()

			filename := tmpFile.Name()
			got, err := Load(&filename)

			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && *got != *tt.want {
				t.Errorf("Load() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &Config{
				Floors:   3,
				Monsters: 2,
				Duration: 1,
			},
			wantErr: false,
		},
		{
			name: "too few floors",
			cfg: &Config{
				Floors:   1,
				Monsters: 2,
				Duration: 1,
			},
			wantErr: true,
		},
		{
			name: "too few monsters",
			cfg: &Config{
				Floors:   2,
				Monsters: 0,
				Duration: 1,
			},
			wantErr: true,
		},
		{
			name: "too few duration",
			cfg: &Config{
				Floors:   2,
				Monsters: 1,
				Duration: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetCloseTime(t *testing.T) {
	cfg := &Config{Duration: 2}
	openTime, _ := time.Parse("15:04:05", "14:00:00")

	closeTime := cfg.GetCloseTime(openTime)
	expected := "16:00:00"

	if closeTime.Format("15:04:05") != expected {
		t.Errorf("GetCloseTime() = %v, want %v", closeTime.Format("15:04:05"), expected)
	}
}
