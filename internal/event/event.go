package event

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"main.go/internal/utils"
)

const (
	EventRegister = iota + 1
	EventEnter
	EventKillMonster
	EventNextFloor
	EventPrevFloor
	EventBossFloor
	EventKillBoss
	EventLeaveDungeon
	EventCannotContinue
	EventRestoreHealth
	EventDamage

	EventDisqualified = iota + 31
	EventDead
	EventImpossibleMove
)

type Event struct {
	Time        time.Time
	PlayerID    int
	ExtraParams string
	Type        int
}

func ParseEvents(eventsPath *string) (*[]Event, error) {
	eventsFile, err := os.Open(*eventsPath)
	if err != nil {
		return &[]Event{}, fmt.Errorf("open events file: %w", err)
	}
	defer eventsFile.Close()

	var events []Event
	scanner := bufio.NewScanner(eventsFile)

	for scanner.Scan() {
		line := scanner.Text()
		event, err := parseEvent(line)
		if err != nil {
			continue
		}
		events = append(events, event)
	}

	return &events, nil
}

func parseEvent(line string) (Event, error) {
	line = strings.TrimSpace(line)
	args := strings.Split(line, " ")
	if len(args) < 3 {
		return Event{}, fmt.Errorf("event has not enough params")
	}

	eventTime, err := utils.ParseTime(strings.Trim(args[0], "[]"))
	if err != nil {
		return Event{}, fmt.Errorf("invalid event time: %w", err)
	}

	playerID, err := strconv.Atoi(args[1])
	if err != nil {
		return Event{}, fmt.Errorf("invalid playerID: %w", err)
	}

	eventType, err := strconv.Atoi(args[2])
	if err != nil {
		return Event{}, fmt.Errorf("invalid eventType: %w", err)
	}

	extraParams := strings.Join(args[3:], " ")

	return Event{
		Time:        eventTime,
		PlayerID:    playerID,
		Type:        eventType,
		ExtraParams: extraParams,
	}, nil
}
