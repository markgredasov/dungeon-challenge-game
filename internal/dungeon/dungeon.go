package dungeon

import (
	"fmt"
	"strconv"
	"time"

	"main.go/internal/config"
	"main.go/internal/event"
	"main.go/internal/state"
	"main.go/internal/utils"
)

type Dungeon struct {
	cfg       *config.Config
	openTime  time.Time
	closeTime time.Time
	players   map[int]*state.State
	output    []string
}

func NewDungeon(cfg *config.Config) (*Dungeon, error) {
	openTime, err := utils.ParseTime(cfg.OpenAt)
	if err != nil {
		return &Dungeon{}, fmt.Errorf("parsing open time: %w", err)
	}

	closeTime := cfg.GetCloseTime(openTime)

	return &Dungeon{
		cfg:       cfg,
		openTime:  openTime,
		closeTime: closeTime,
		players:   make(map[int]*state.State),
	}, nil
}

func (d *Dungeon) HandleEvent(e *event.Event) error {
	if e.Time.After(d.closeTime) {
		return fmt.Errorf("event after dungeon closed: %v", e)
	}

	p := d.addOrCheckPlayer(e.PlayerID)

	switch e.Type {
	case event.EventRegister:
		d.registerPlayer(p, e)
	case event.EventEnter:
		d.enterTheDungeon(p, e)
	case event.EventKillMonster:
		d.killMonster(p, e)
	case event.EventNextFloor:
		d.toNextFloor(p, e)
	case event.EventPrevFloor:
		d.toPrevFloor(p, e)
	case event.EventBossFloor:
		d.toBossFloor(p, e)
	case event.EventKillBoss:
		d.killTheBoss(p, e)
	case event.EventLeaveDungeon:
		d.leaveTheDungeon(p, e)
	case event.EventCannotContinue:
		d.cannotContinue(p, e)
	case event.EventRestoreHealth:
		d.restoreHealth(p, e)
	case event.EventDamage:
		d.getDamage(p, e)
	default:
		return fmt.Errorf("wrong event type: %v", e)
	}

	return nil
}

func (d *Dungeon) addOrCheckPlayer(playerID int) *state.State {
	if p, ok := d.players[playerID]; ok {
		return p
	}

	p := state.NewPlayer(playerID)
	d.players[playerID] = p
	return p
}

func (d *Dungeon) addOutput(e *event.Event, text string) {
	t := utils.TimeToStr(e.Time)
	d.output = append(d.output, fmt.Sprintf("[%s] %s", t, text))
}

func (d *Dungeon) registerPlayer(p *state.State, e *event.Event) {
	if p.IsRegistered {
		return
	}
	p.IsRegistered = true
	d.addOutput(e, fmt.Sprintf("Player [%d] registered", p.ID))
}

func (d *Dungeon) enterTheDungeon(p *state.State, e *event.Event) {
	if !p.IsRegistered {
		p.IsDisqualified = true
		d.addOutput(e, fmt.Sprintf("Player [%d] is disqualified", p.ID))
		return
	}

	if e.Time.Before(d.openTime) {
		p.IsDisqualified = true
		d.addOutput(e, fmt.Sprintf("Player [%d] is disqualified", p.ID))
		return
	}

	if !p.InDungeon {
		p.InDungeon = true
		p.DungeonEnterTime = e.Time
		p.FloorEnterTime = e.Time
		d.addOutput(e, fmt.Sprintf("Player [%d] entered the dungeon", p.ID))
	}
}

func (d *Dungeon) killMonster(p *state.State, e *event.Event) {
	if !p.CanAct() {
		return
	}

	if p.CurrentFloor >= d.cfg.Floors {
		return
	}

	maxMonsters := d.cfg.Monsters
	if p.MonstersKilled[p.CurrentFloor] < maxMonsters {
		p.KillMonster(p.CurrentFloor)
		d.addOutput(e, fmt.Sprintf("Player [%d] killed the monster", p.ID))

		if p.MonstersKilled[p.CurrentFloor] == maxMonsters && !p.ClearedFloors[p.CurrentFloor] {
			clearTime := e.Time.Sub(p.FloorEnterTime)
			p.FloorClearTimeSlice = append(p.FloorClearTimeSlice, clearTime)
			p.ClearedFloors[p.CurrentFloor] = true
		}
	} else {
		d.addOutput(e, fmt.Sprintf("Player [%d] makes imposible move [%d]", p.ID, e.Type))
	}
}

func (d *Dungeon) toNextFloor(p *state.State, e *event.Event) {
	if !p.CanAct() {
		return
	}

	if p.CurrentFloor >= d.cfg.Floors {
		d.addOutput(e, fmt.Sprintf("Player [%d] makes imposible move [%d]", p.ID, e.Type))
		return
	}

	if !p.ClearedFloors[p.CurrentFloor] {
		d.addOutput(e, fmt.Sprintf("Player [%d] makes imposible move [%d]", p.ID, e.Type))
		return
	}

	p.CurrentFloor++
	p.FloorEnterTime = e.Time
	d.addOutput(e, fmt.Sprintf("Player [%d] went to the next floor", p.ID))
}

func (d *Dungeon) toPrevFloor(p *state.State, e *event.Event) {
	if !p.CanAct() {
		return
	}

	if p.CurrentFloor > 1 {
		p.CurrentFloor--
		d.addOutput(e, fmt.Sprintf("Player [%d] went to the previous floor", p.ID))
	} else {
		d.addOutput(e, fmt.Sprintf("Player [%d] makes imposible move [%d]", p.ID, e.Type))
	}
}

func (d *Dungeon) toBossFloor(p *state.State, e *event.Event) {
	if !p.CanAct() {
		return
	}

	if p.CurrentFloor != d.cfg.Floors {
		d.addOutput(e, fmt.Sprintf("Player [%d] makes imposible move [%d]", p.ID, e.Type))
		return
	}

	d.addOutput(e, fmt.Sprintf("Player [%d] entered the boss's floor", p.ID))
}

func (d *Dungeon) killTheBoss(p *state.State, e *event.Event) {
	if !p.CanAct() {
		return
	}

	if p.CurrentFloor != d.cfg.Floors {
		d.addOutput(e, fmt.Sprintf("Player [%d] makes imposible move [%d]", p.ID, e.Type))
		return
	}

	p.KillTheBossTime = e.Time
	d.addOutput(e, fmt.Sprintf("Player [%d] killed the boss", p.ID))
}

func (d *Dungeon) leaveTheDungeon(p *state.State, e *event.Event) {
	if !p.CanAct() {
		return
	}

	p.InDungeon = false
	p.DungeonExitTime = e.Time
	d.addOutput(e, fmt.Sprintf("Player [%d] left the dungeon", p.ID))
}

func (d *Dungeon) cannotContinue(p *state.State, e *event.Event) {
	p.IsDisqualified = true
	p.InDungeon = false
	p.DungeonExitTime = e.Time
	d.addOutput(e, fmt.Sprintf("Player [%d] cannot continue due to [%s]", p.ID, e.ExtraParams))
}

func (d *Dungeon) restoreHealth(p *state.State, e *event.Event) {
	if !p.CanAct() {
		return
	}

	amount, err := strconv.Atoi(e.ExtraParams)
	if err != nil {
		return
	}

	p.Heal(amount)
	d.addOutput(e, fmt.Sprintf("Player [%d] has restored [%s] of health", p.ID, e.ExtraParams))
}

func (d *Dungeon) getDamage(p *state.State, e *event.Event) {
	if !p.CanAct() {
		return
	}

	amount, err := strconv.Atoi(e.ExtraParams)
	if err != nil {
		return
	}

	p.GetDamage(amount)
	d.addOutput(e, fmt.Sprintf("Player [%d] received [%s] of damage", p.ID, e.ExtraParams))

	if !p.IsAlive {
		p.DungeonExitTime = e.Time
		d.playerDead(p, e)
	}
}

func (d *Dungeon) playerDead(p *state.State, e *event.Event) {
	d.addOutput(e, fmt.Sprintf("Player [%d] is dead", p.ID))
}

func (d *Dungeon) Finalize() {
	for _, p := range d.players {
		if p.InDungeon {
			p.DungeonExitTime = d.closeTime
			p.InDungeon = false
		}
	}
}

func (d *Dungeon) PrintOutput() {
	for _, line := range d.output {
		fmt.Println(line)
	}
}

func (d *Dungeon) GetPlayers() map[int]*state.State {
	return d.players
}

func (d *Dungeon) GetConfig() *config.Config {
	return d.cfg
}

func (d *Dungeon) CloseTime() time.Time {
	return d.closeTime
}
