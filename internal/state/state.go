package state

import "time"

type State struct {
	ID                  int
	Health              int
	CurrentFloor        int
	IsRegistered        bool
	InDungeon           bool
	IsAlive             bool
	IsDisqualified      bool
	MonstersKilled      map[int]int
	ClearedFloors       map[int]bool
	DungeonEnterTime    time.Time
	FloorEnterTime      time.Time
	KillTheBossTime     time.Time
	DungeonExitTime     time.Time
	FloorClearTimeSlice []time.Duration
}

func NewPlayer(id int) *State {
	return &State{
		ID:             id,
		Health:         100,
		InDungeon:      false,
		IsAlive:        true,
		CurrentFloor:   1,
		MonstersKilled: make(map[int]int),
		ClearedFloors:  make(map[int]bool),
	}
}

func (p *State) CanAct() bool {
	return !p.IsDisqualified && p.InDungeon && p.IsAlive
}

func (p *State) KillMonster(floor int) {
	p.MonstersKilled[floor]++
}

func (p *State) Heal(amount int) {
	result := p.Health + amount
	if result > 100 {
		p.Health = 100
	} else {
		p.Health = result
	}
}

func (p *State) GetDamage(amount int) {
	result := p.Health - amount
	if result <= 0 {
		p.IsAlive = false
		p.Health = 0
		p.InDungeon = false
	} else {
		p.Health = result
	}
}

func (p *State) AreAllFloorsCleared(nFloors int, nMonsters int) bool {
	for floor := 1; floor < nFloors; floor++ {
		if p.MonstersKilled[floor] < nMonsters {
			return false
		}
	}
	return true
}
