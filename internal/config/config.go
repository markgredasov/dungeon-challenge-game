package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Floors   int    `json:"Floors"`
	Monsters int    `json:"Monsters"`
	OpenAt   string `json:"OpenAt"`
	Duration int    `json:"Duration"`
}

func Load(configFile *string) (*Config, error) {
	data, err := os.ReadFile(*configFile)
	if err != nil {
		return &Config{}, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return &Config{}, fmt.Errorf("unmarshalling config file: %w", err)
	}

	if err := Validate(&cfg); err != nil {
		return &Config{}, fmt.Errorf("validating config: %w", err)
	}

	return &cfg, nil
}

func Validate(cfg *Config) error {
	if cfg.Floors < 2 {
		return fmt.Errorf("too few floors")
	}
	if cfg.Monsters < 1 {
		return fmt.Errorf("too few monsters")
	}
	if cfg.Duration < 1 {
		return fmt.Errorf("too few duration")
	}
	return nil
}

func (cfg *Config) GetCloseTime(openTime time.Time) time.Time {
	return openTime.Add(time.Duration(cfg.Duration) * time.Hour)
}
