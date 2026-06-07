package main

import (
	"flag"
	"fmt"
	"log"

	"main.go/internal/config"
	"main.go/internal/dungeon"
	"main.go/internal/event"
	"main.go/internal/report"
)

func main() {
	configFile := flag.String("config", "config.json", "configuration")
	eventsFile := flag.String("events", "events", "events of players")
	flag.Parse()

	cfg, err := config.Load(configFile)
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	d, err := dungeon.NewDungeon(cfg)
	if err != nil {
		log.Fatalf("building dungeon: %v", err)
	}

	events, err := event.ParseEvents(eventsFile)
	if err != nil {
		log.Fatalf("parsing events: %v", err)
	}

	for _, e := range *events {
		if err := d.HandleEvent(&e); err != nil {
			log.Printf("handle event: %v\n", err)
		}
	}

	d.Finalize()
	d.PrintOutput()
	fmt.Println(report.Generate(d))
}
