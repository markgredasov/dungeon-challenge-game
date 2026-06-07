package dungeon

import (
	"testing"
	"time"

	"main.go/internal/config"
	"main.go/internal/event"
)

func setupTestDungeon() *Dungeon {
	cfg := &config.Config{
		Floors:   3,
		Monsters: 2,
		OpenAt:   "14:00:00",
		Duration: 2,
	}

	d, _ := NewDungeon(cfg)
	return d
}

func makeEvent(eventType int, playerID int, timeStr string, extraParams ...string) *event.Event {
	t, _ := time.Parse("15:04:05", timeStr)
	extra := ""
	if len(extraParams) > 0 {
		extra = extraParams[0]
	}
	return &event.Event{
		Time:        t,
		PlayerID:    playerID,
		Type:        eventType,
		ExtraParams: extra,
	}
}

func TestRegisterPlayer(t *testing.T) {
	d := setupTestDungeon()
	e := makeEvent(event.EventRegister, 1, "14:00:00")

	if err := d.HandleEvent(e); err != nil {
		t.Fatalf("HandleEvent() error = %v", err)
	}

	p := d.players[1]
	if !p.IsRegistered {
		t.Error("Player should be registered")
	}
}

func TestEnterDungeonUnregistered(t *testing.T) {
	d := setupTestDungeon()
	e := makeEvent(event.EventEnter, 1, "14:01:00")

	d.HandleEvent(e)

	p := d.players[1]
	if !p.IsDisqualified {
		t.Error("Unregistered player should be disqualified")
	}

	expected := "[14:01:00] Player [1] is disqualified"
	if d.output[0] != expected {
		t.Errorf("Output = %v, want %v", d.output[0], expected)
	}
}

func TestEnterDungeonBeforeOpen(t *testing.T) {
	d := setupTestDungeon()

	regEvent := makeEvent(event.EventRegister, 1, "13:00:00")
	d.HandleEvent(regEvent)

	enterEvent := makeEvent(event.EventEnter, 1, "13:59:00")
	d.HandleEvent(enterEvent)

	p := d.players[1]
	if !p.IsDisqualified {
		t.Error("Player should be disqualified for entering before opening")
	}
}

func TestFullGameFlow(t *testing.T) {
	d := setupTestDungeon()

	events := []*event.Event{
		makeEvent(event.EventRegister, 1, "14:00:00"),
		makeEvent(event.EventEnter, 1, "14:00:01"),
		makeEvent(event.EventKillMonster, 1, "14:00:02"),
		makeEvent(event.EventKillMonster, 1, "14:00:03"),
		makeEvent(event.EventNextFloor, 1, "14:00:04"),
		makeEvent(event.EventKillMonster, 1, "14:00:05"),
		makeEvent(event.EventKillMonster, 1, "14:00:06"),
		makeEvent(event.EventNextFloor, 1, "14:00:07"),
		makeEvent(event.EventBossFloor, 1, "14:00:08"),
		makeEvent(event.EventKillBoss, 1, "14:00:09"),
		makeEvent(event.EventLeaveDungeon, 1, "14:00:10"),
	}

	for _, e := range events {
		if err := d.HandleEvent(e); err != nil {
			t.Fatalf("HandleEvent() error = %v at event %v", err, e)
		}
	}

	p := d.players[1]
	if !p.AreAllFloorsCleared(d.cfg.Floors, d.cfg.Monsters) {
		t.Error("All floors should be cleared")
	}

	if p.KillTheBossTime.IsZero() {
		t.Error("Boss should be killed")
	}
}

func TestDamageAndDeath(t *testing.T) {
	d := setupTestDungeon()

	events := []*event.Event{
		makeEvent(event.EventRegister, 1, "14:00:00"),
		makeEvent(event.EventEnter, 1, "14:00:01"),
		makeEvent(event.EventDamage, 1, "14:00:02", "50"),
		makeEvent(event.EventDamage, 1, "14:00:03", "60"),
	}

	for _, e := range events {
		d.HandleEvent(e)
	}

	p := d.players[1]
	if p.IsAlive {
		t.Error("Player should be dead")
	}
	if p.Health != 0 {
		t.Errorf("Health = %v, want 0", p.Health)
	}
}

func TestImpossibleMove(t *testing.T) {
	tests := []struct {
		name   string
		events []*event.Event
		want   string
	}{
		{
			name: "move to next floor without clearing",
			events: []*event.Event{
				makeEvent(event.EventRegister, 1, "14:00:00"),
				makeEvent(event.EventEnter, 1, "14:00:01"),
				makeEvent(event.EventNextFloor, 1, "14:00:02"),
			},
			want: "[14:00:02] Player [1] makes imposible move [4]",
		},
		{
			name: "move to prev floor from floor 1",
			events: []*event.Event{
				makeEvent(event.EventRegister, 1, "14:00:00"),
				makeEvent(event.EventEnter, 1, "14:00:01"),
				makeEvent(event.EventPrevFloor, 1, "14:00:02"),
			},
			want: "[14:00:02] Player [1] makes imposible move [5]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := setupTestDungeon()

			for _, e := range tt.events {
				d.HandleEvent(e)
			}

			found := false
			for _, out := range d.output {
				if out == tt.want {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("Expected output '%s' not found in %v", tt.want, d.output)
			}
		})
	}
}

func TestRestoreHealth(t *testing.T) {
	d := setupTestDungeon()

	events := []*event.Event{
		makeEvent(event.EventRegister, 1, "14:00:00"),
		makeEvent(event.EventEnter, 1, "14:00:01"),
		makeEvent(event.EventDamage, 1, "14:00:02", "50"),
		makeEvent(event.EventRestoreHealth, 1, "14:00:03", "30"),
	}

	for _, e := range events {
		d.HandleEvent(e)
	}

	p := d.players[1]
	if p.Health != 80 {
		t.Errorf("Health = %v, want 80", p.Health)
	}
}

func TestCannotContinue(t *testing.T) {
	d := setupTestDungeon()

	events := []*event.Event{
		makeEvent(event.EventRegister, 1, "14:00:00"),
		makeEvent(event.EventEnter, 1, "14:00:01"),
		makeEvent(event.EventCannotContinue, 1, "14:00:05", "is tired"),
	}

	for _, e := range events {
		d.HandleEvent(e)
	}

	p := d.players[1]
	if !p.IsDisqualified {
		t.Error("Player should be disqualified")
	}
}

func TestDungeonFinalize(t *testing.T) {
	d := setupTestDungeon()

	events := []*event.Event{
		makeEvent(event.EventRegister, 1, "14:00:00"),
		makeEvent(event.EventEnter, 1, "15:30:00"),
	}

	for _, e := range events {
		d.HandleEvent(e)
	}

	d.Finalize()

	p := d.players[1]
	if p.InDungeon {
		t.Error("Player should not be in dungeon after finalize")
	}
	if p.DungeonExitTime.IsZero() {
		t.Error("Exit time should be set")
	}
}
