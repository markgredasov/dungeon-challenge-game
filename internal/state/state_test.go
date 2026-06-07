package state

import (
	"testing"
)

func TestNewPlayer(t *testing.T) {
	p := NewPlayer(1)

	if p.ID != 1 {
		t.Errorf("NewPlayer() ID = %v, want 1", p.ID)
	}
	if p.Health != 100 {
		t.Errorf("NewPlayer() Health = %v, want 100", p.Health)
	}
	if !p.IsAlive {
		t.Errorf("NewPlayer() IsAlive = %v, want true", p.IsAlive)
	}
	if p.CurrentFloor != 1 {
		t.Errorf("NewPlayer() CurrentFloor = %v, want 1", p.CurrentFloor)
	}
}

func TestCanAct(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*State)
		wantCanAct bool
	}{
		{
			name: "normal player can act",
			setup: func(s *State) {
				s.InDungeon = true
			},
			wantCanAct: true,
		},
		{
			name:       "player not in dungeon",
			setup:      func(s *State) {},
			wantCanAct: false,
		},
		{
			name: "dead player",
			setup: func(s *State) {
				s.InDungeon = true
				s.IsAlive = false
			},
			wantCanAct: false,
		},
		{
			name: "disqualified player",
			setup: func(s *State) {
				s.InDungeon = true
				s.IsDisqualified = true
			},
			wantCanAct: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPlayer(1)
			tt.setup(p)

			if got := p.CanAct(); got != tt.wantCanAct {
				t.Errorf("CanAct() = %v want %v", got, tt.wantCanAct)
			}
		})
	}
}

func TestHeal(t *testing.T) {
	tests := []struct {
		name       string
		initialHP  int
		healAmount int
		expectedHP int
	}{
		{
			name:       "normal heal",
			initialHP:  50,
			healAmount: 30,
			expectedHP: 80,
		},
		{
			name:       "heal above max",
			initialHP:  80,
			healAmount: 50,
			expectedHP: 100,
		},
		{
			name:       "heal to exactly max",
			initialHP:  70,
			healAmount: 30,
			expectedHP: 100,
		},
		{
			name:       "heal from max",
			initialHP:  100,
			healAmount: 10,
			expectedHP: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPlayer(1)
			p.Health = tt.initialHP
			p.Heal(tt.healAmount)

			if p.Health != tt.expectedHP {
				t.Errorf("Heal() Health = %v want %v", p.Health, tt.expectedHP)
			}
		})
	}
}

func TestGetDamage(t *testing.T) {
	tests := []struct {
		name              string
		initialHP         int
		damage            int
		expectedHP        int
		expectedAlive     bool
		expectedInDungeon bool
	}{
		{
			name:              "normal damage",
			initialHP:         100,
			damage:            30,
			expectedHP:        70,
			expectedAlive:     true,
			expectedInDungeon: true,
		},
		{
			name:              "fatal damage",
			initialHP:         50,
			damage:            60,
			expectedHP:        0,
			expectedAlive:     false,
			expectedInDungeon: false,
		},
		{
			name:              "exactly zero HP",
			initialHP:         50,
			damage:            50,
			expectedHP:        0,
			expectedAlive:     false,
			expectedInDungeon: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPlayer(1)
			p.Health = tt.initialHP
			p.InDungeon = true
			p.GetDamage(tt.damage)

			if p.Health != tt.expectedHP {
				t.Errorf("GetDamage() Health = %v, want %v", p.Health, tt.expectedHP)
			}
			if p.IsAlive != tt.expectedAlive {
				t.Errorf("GetDamage() IsAlive = %v, want %v", p.IsAlive, tt.expectedAlive)
			}
			if p.InDungeon != tt.expectedInDungeon {
				t.Errorf("GetDamage() InDungeon = %v, want %v", p.InDungeon, tt.expectedInDungeon)
			}
		})
	}
}

func TestAreAllFloorsCleared(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*State)
		nFloors   int
		nMonsters int
		want      bool
	}{
		{
			name: "all floors cleared",
			setup: func(s *State) {
				s.MonstersKilled[1] = 2
				s.MonstersKilled[2] = 2
			},
			nFloors:   3,
			nMonsters: 2,
			want:      true,
		},
		{
			name: "not all floors cleared",
			setup: func(s *State) {
				s.MonstersKilled[1] = 2
				s.MonstersKilled[2] = 1
			},
			nFloors:   3,
			nMonsters: 2,
			want:      false,
		},
		{
			name:      "empty floors",
			setup:     func(s *State) {},
			nFloors:   2,
			nMonsters: 1,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPlayer(1)
			tt.setup(p)

			if got := p.AreAllFloorsCleared(tt.nFloors, tt.nMonsters); got != tt.want {
				t.Errorf("AreAllFloorsCleared() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKillMonster(t *testing.T) {
	p := NewPlayer(1)

	p.KillMonster(1)
	if p.MonstersKilled[1] != 1 {
		t.Errorf("KillMonster() = %v, want 1", p.MonstersKilled[1])
	}

	p.KillMonster(1)
	if p.MonstersKilled[1] != 2 {
		t.Errorf("KillMonster() = %v, want 2", p.MonstersKilled[1])
	}
}
