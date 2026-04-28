package game

import (
	"fmt"
	"math/rand"
	"strings"
)

const threatTickInterval = 3.5

type HitZone string

const (
	ZoneHead HitZone = "head"
	ZoneCore HitZone = "core"
	ZoneLegs HitZone = "legs"
)

type TrainThreatType string

const (
	ThreatFire   TrainThreatType = "fire"
	ThreatBreach TrainThreatType = "breach"
	ThreatShort  TrainThreatType = "short"
)

type TrainThreat struct {
	ID            int
	RoomID        int
	Type          TrainThreatType
	Severity      int
	TickRemaining float32
}

type WaveState struct {
	Number               int
	TimeRemaining        float32
	KillsRequired        int
	KillsDone            int
	ThreatSpawnRemaining float32
	Active               bool
	Success              bool
	Failed               bool
}

func parseHitZone(raw string) (HitZone, bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case string(ZoneHead):
		return ZoneHead, true
	case string(ZoneCore), "hull":
		return ZoneCore, true
	case string(ZoneLegs):
		return ZoneLegs, true
	default:
		return "", false
	}
}

func (g *Game) startWave(number int) {
	if number <= 0 {
		number = 1
	}

	timeBudget := float32((50 + number*15) * 3)
	killsRequired := 2 + number
	if killsRequired > 8 {
		killsRequired = 8
	}

	g.Wave = WaveState{
		Number:               number,
		TimeRemaining:        timeBudget,
		KillsRequired:        killsRequired,
		KillsDone:            0,
		ThreatSpawnRemaining: 5,
		Active:               true,
		Success:              false,
		Failed:               false,
	}
	g.ActiveThreats = g.ActiveThreats[:0]
	g.spawnEnemyForWave()
	g.SetCombatStatus(
		fmt.Sprintf("Wave %d started. Eliminate %d hostiles in %.0fs.", g.Wave.Number, g.Wave.KillsRequired, g.Wave.TimeRemaining),
		"Use 'target zone <head|core|legs>' then 'fire <weapon>'.",
	)
}

func (g *Game) spawnEnemyForWave() {
	hp := 14 + g.Wave.Number*4
	attack := 1 + (g.Wave.Number-1)/2
	if attack < 1 {
		attack = 1
	}
	g.Enemy = NewBasicEnemy(fmt.Sprintf("Raider %d", g.Wave.KillsDone+1), hp, attack)
}

func (g *Game) updateWaveState(deltaSeconds float32) {
	if !g.Wave.Active || deltaSeconds <= 0 {
		return
	}

	g.Wave.TimeRemaining -= deltaSeconds
	if g.Wave.TimeRemaining <= 0 {
		g.Wave.TimeRemaining = 0
		g.Wave.Active = false
		g.Wave.Failed = true
		g.SetCombatStatus(
			fmt.Sprintf("Wave %d failed. Time expired.", g.Wave.Number),
			"Tip: prioritize repairs to keep systems online.",
		)
		return
	}

	if g.Train.Health <= 0 {
		g.Wave.Active = false
		g.Wave.Failed = true
		g.SetCombatStatus("Wave failed. Train hull collapsed.")
		return
	}

	g.Wave.ThreatSpawnRemaining -= deltaSeconds
	if g.Wave.ThreatSpawnRemaining <= 0 {
		g.spawnThreat()
		g.Wave.ThreatSpawnRemaining = float32(8 - g.Wave.Number)
		if g.Wave.ThreatSpawnRemaining < 4 {
			g.Wave.ThreatSpawnRemaining = 4
		}
	}

	g.tickThreats(deltaSeconds)
}

func (g *Game) tickThreats(deltaSeconds float32) {
	if len(g.ActiveThreats) == 0 {
		return
	}

	for i := range g.ActiveThreats {
		g.ActiveThreats[i].TickRemaining -= deltaSeconds
		if g.ActiveThreats[i].TickRemaining > 0 {
			continue
		}

		g.ActiveThreats[i].TickRemaining += threatTickInterval
		g.applyThreatTick(g.ActiveThreats[i])
	}
}

func (g *Game) applyThreatTick(threat TrainThreat) {
	room, ok := g.Train.GetRoom(threat.RoomID)
	if !ok {
		return
	}

	roomDamage := threat.Severity
	hullDamage := 0

	switch threat.Type {
	case ThreatBreach:
		hullDamage = threat.Severity
	case ThreatFire:
		if threat.Severity >= 2 {
			hullDamage = 1
		}
	case ThreatShort:
		if room.System.Type == SystemWeapons || room.System.Type == SystemShields || room.System.Type == SystemEngines {
			roomDamage++
		}
	}

	g.Train.ApplyRoomDamage(threat.RoomID, roomDamage)
	if hullDamage > 0 {
		g.Train.ApplyHullDamage(hullDamage)
	}
}

func (g *Game) spawnThreat() {
	if len(g.Train.Rooms) == 0 || len(g.ActiveThreats) >= 4 {
		return
	}

	candidateRooms := make([]int, 0, len(g.Train.Rooms))
	for roomIdx := range g.Train.Rooms {
		if !g.hasThreatInRoom(roomIdx) {
			candidateRooms = append(candidateRooms, roomIdx)
		}
	}

	if len(candidateRooms) == 0 {
		return
	}

	roomID := candidateRooms[rand.Intn(len(candidateRooms))]
	threatTypes := []TrainThreatType{ThreatFire, ThreatBreach, ThreatShort}
	threatType := threatTypes[rand.Intn(len(threatTypes))]
	severity := 1 + rand.Intn(2)

	g.nextThreatID++
	g.ActiveThreats = append(g.ActiveThreats, TrainThreat{
		ID:            g.nextThreatID,
		RoomID:        roomID,
		Type:          threatType,
		Severity:      severity,
		TickRemaining: threatTickInterval,
	})

	roomRune := string(g.Train.Rooms[roomID].GetRune())
	g.SetCombatStatus(
		fmt.Sprintf("New threat in room %s: %s (severity %d).", roomRune, strings.ToUpper(string(threatType)), severity),
		"Move crew and run 'repair <room>' to stabilize the car.",
	)
}

func (g *Game) hasThreatInRoom(roomID int) bool {
	for _, threat := range g.ActiveThreats {
		if threat.RoomID == roomID {
			return true
		}
	}

	return false
}

func (g *Game) resolveThreatInRoom(roomID int, crewSupport int) bool {
	if crewSupport <= 0 {
		crewSupport = 1
	}

	for i := range g.ActiveThreats {
		if g.ActiveThreats[i].RoomID != roomID {
			continue
		}

		g.ActiveThreats[i].Severity -= crewSupport

		if g.ActiveThreats[i].Severity <= 0 {
			g.ActiveThreats = append(g.ActiveThreats[:i], g.ActiveThreats[i+1:]...)
		}
		return true
	}

	return false
}

func (g *Game) crewInRoom(roomID int) bool {
	for _, character := range g.Train.Characters {
		if character.Pos.RoomId == roomID && !character.IsMoving {
			return true
		}
	}

	return false
}

func (g *Game) registerEnemyKill() {
	if !g.Wave.Active {
		return
	}

	g.Wave.KillsDone++
	if g.Wave.KillsDone >= g.Wave.KillsRequired {
		g.Wave.Active = false
		g.Wave.Success = true
		g.Enemy = nil
		g.SetCombatStatus(
			fmt.Sprintf("Wave %d complete.", g.Wave.Number),
			"All hostiles neutralized.",
		)
		return
	}

	g.spawnEnemyForWave()
}

func (g *Game) WaveSummaryText() string {
	status := "ACTIVE"
	if g.Wave.Success {
		status = "COMPLETE"
	} else if g.Wave.Failed {
		status = "FAILED"
	}

	return fmt.Sprintf(
		"Wave %d [%s]  Time %.0fs  Kills %d/%d",
		g.Wave.Number,
		status,
		g.Wave.TimeRemaining,
		g.Wave.KillsDone,
		g.Wave.KillsRequired,
	)
}

func (g *Game) ThreatSummaryText() string {
	if len(g.ActiveThreats) == 0 {
		return "none"
	}

	damageByRoom := make(map[int]int, len(g.ActiveThreats))
	for _, threat := range g.ActiveThreats {
		damageByRoom[threat.RoomID] += threat.Severity
	}

	parts := make([]string, 0, len(damageByRoom))
	for roomID := range g.Train.Rooms {
		damage := damageByRoom[roomID]
		if damage <= 0 {
			continue
		}

		parts = append(parts, fmt.Sprintf("%s -> %d damage", g.cartLabel(roomID), damage))
	}

	if len(parts) == 0 {
		return "none"
	}

	return strings.Join(parts, "  |  ")
}

func weaponPreferredZone(weapon Weapon) HitZone {
	switch weapon.Type {
	case WeaponMissile:
		return ZoneHead
	case WeaponCannon:
		return ZoneCore
	default:
		return ZoneLegs
	}
}

func weaponDamageForZone(weapon Weapon, zone HitZone) int {
	damage := weapon.Damage
	preferred := weaponPreferredZone(weapon)

	if zone == preferred {
		damage += 2
	} else if zone == ZoneLegs {
		damage--
	}

	if damage < 1 {
		damage = 1
	}

	return damage
}
