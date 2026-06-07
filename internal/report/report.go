package report

import (
	"fmt"
	"strings"
	"time"

	"main.go/internal/dungeon"
	"main.go/internal/state"
	"main.go/internal/utils"
)

type PlayerReport struct {
	Status       string
	PlayerID     int
	DungeonTime  time.Duration
	AvgClearTime time.Duration
	BossTime     time.Duration
	Health       int
}

func Generate(d *dungeon.Dungeon) string {
	var reports []PlayerReport

	for _, p := range d.GetPlayers() {
		report := calculatePlayerReport(p, d)
		reports = append(reports, report)
	}

	return formatReports(reports)
}

func calculatePlayerReport(p *state.State, d *dungeon.Dungeon) PlayerReport {
	report := PlayerReport{
		PlayerID: p.ID,
		Health:   p.Health,
	}

	if !p.IsRegistered || p.IsDisqualified {
		report.Status = "DISQUAL"
		return report
	}

	allMonstersCleared := p.AreAllFloorsCleared(d.GetConfig().Floors, d.GetConfig().Monsters)
	bossKilled := !p.KillTheBossTime.IsZero()

	if allMonstersCleared && bossKilled {
		report.Status = "SUCCESS"
	} else {
		report.Status = "FAIL"
	}

	if !p.DungeonExitTime.IsZero() {
		report.DungeonTime = p.DungeonExitTime.Sub(p.DungeonEnterTime)
	} else if !p.DungeonEnterTime.IsZero() {
		report.DungeonTime = d.CloseTime().Sub(p.DungeonEnterTime)
	}

	if len(p.FloorClearTimeSlice) > 0 {
		var total time.Duration
		for _, t := range p.FloorClearTimeSlice {
			total += t
		}
		report.AvgClearTime = total / time.Duration(len(p.FloorClearTimeSlice))
	}

	if !p.KillTheBossTime.IsZero() {
		report.BossTime = p.KillTheBossTime.Sub(p.FloorEnterTime)
		if report.BossTime < 0 {
			report.BossTime = 0
		}
	}

	return report
}

func formatReports(reports []PlayerReport) string {
	var sb strings.Builder
	sb.WriteString("\nFinal report:\n")

	for _, r := range reports {
		sb.WriteString(fmt.Sprintf("[%s] %d [%s, %s, %s] HP:%d\n",
			r.Status,
			r.PlayerID,
			utils.FormatDuration(r.DungeonTime),
			utils.FormatDuration(r.AvgClearTime),
			utils.FormatDuration(r.BossTime),
			r.Health,
		))
	}

	return sb.String()
}
