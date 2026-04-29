package game

import (
	"fmt"
	"strings"

	"github.com/KydZombie/armada/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	Train *Train
	Enemy Enemy

	windows  []core.Window[Game]
	Terminal *Terminal

	SelectedCharacterIndex int
	crewSystemTickTimer    float32
	combatStatusLines      []string

	SelectedWeaponIndex   int
	SelectedTargetZone    HitZone
	Wave                  WaveState
	ActiveThreats         []TrainThreat
	nextThreatID          int
	MissionBriefingActive bool
}

func NewGameScreen(gm *core.GameManager) *Game {
	train := NewTrain(100)
	enemy := NewBasicEnemy("Steel Matador", 20, 3)

	gs := &Game{
		Train: train,
		Enemy: enemy,

		windows: []core.Window[Game]{},
		Terminal: &Terminal{
			commandDB: initializeCommands(),
		},

		SelectedCharacterIndex: -1,
		crewSystemTickTimer:    0,
		combatStatusLines:      []string{"No combat actions yet."},
		SelectedTargetZone:     ZoneCore,
		ActiveThreats:          []TrainThreat{},
		MissionBriefingActive:  true,
	}
	gs.startWave(1)

	const windowMargin = 16.0
	const rightColumnInset = 24.0
	const missionExtraInset = 28.0
	const missionWidthScale = 0.7
	const trainMissionGap = 4.0

	terminal := NewTerminalWindow(
		func(gm *core.GameManager) rl.Rectangle {
			return rl.Rectangle{
				X:      64,
				Y:      432,
				Width:  288,
				Height: 256,
			}
		},
		gm,
		gs.Terminal,
	)

	gs.windows = append(gs.windows, terminal)

	gs.windows = append(gs.windows, NewMissionWindow(
		func(gm *core.GameManager) rl.Rectangle {
			return rl.Rectangle{
				X:      928,
				Y:      432,
				Width:  288,
				Height: 240,
			}
		},
		gm,
	))

	gs.windows = append(gs.windows, NewBattleWindow(
		func(gm *core.GameManager) rl.Rectangle {
			return rl.Rectangle{
				X:      float32(gm.ScreenWidth)/2.0 + windowMargin + rightColumnInset,
				Y:      float32(gm.ScreenHeight)/2.0 + windowMargin,
				Width:  float32(gm.ScreenWidth)/2.0 - windowMargin*2 - rightColumnInset,
				Height: float32(gm.ScreenHeight)/2.0 - windowMargin*2,
			}
		},
		gm,
	))

	gs.windows = append(gs.windows, NewTrainWindow(
		func(gm *core.GameManager) rl.Rectangle {
			missionWidth := (float32(gm.ScreenWidth)/2.0 - windowMargin*2 - rightColumnInset - missionExtraInset) * missionWidthScale
			return rl.Rectangle{
				X:      210,
				Y:      10,
				Width:  float32(gm.ScreenWidth) - windowMargin*2 - missionWidth - trainMissionGap,
				Height: float32(gm.ScreenHeight)/2.0 - windowMargin,
			}
		},
		gm,
	))

	return gs
}

func (g *Game) ResizeScreen(gm *core.GameManager) {
	g.UpdateWindowSizes(gm)
}

func (g *Game) UpdateScreen(gm *core.GameManager) {
	inputCaptured := false
	deltaSeconds := rl.GetFrameTime()

	for _, window := range g.windows {
		// If a window captures the input, other windows should not read any input
		if !inputCaptured && window.HandleInput(gm, g) {
			inputCaptured = true
		}

		window.UpdateWindow(gm, g)
	}

	if g.isGameOverModalActive() {
		return
	}
	if g.isMissionBriefingActive() {
		return
	}

	if g.Train.WeaponsOperational() {
		g.Train.AdvanceWeaponCooldowns(deltaSeconds)
	}
	g.Train.UpdateCombatState(deltaSeconds)
	if g.Enemy != nil {
		g.Enemy.UpdateCombatState(deltaSeconds)
	}
	g.updateCrewSystems()
	g.updateWaveState(deltaSeconds)

	// Update character animations
	g.Train.UpdateCharacterAnimations(rl.GetFrameTime())
}

func (g *Game) DrawScreen(gm *core.GameManager) {
	rl.ClearBackground(rl.Black)

	for _, window := range g.windows {
		window.DrawWindow(gm, g)
	}
}

func (g *Game) DrawScreenUI(gm *core.GameManager) {
	if g.isGameOverModalActive() {
		for _, window := range g.windows {
			if _, ok := window.(*MissionWindow); ok {
				window.DrawWindowUI(gm, g)
			}
		}
		return
	}
	if g.isMissionBriefingActive() {
		for _, window := range g.windows {
			if _, ok := window.(*MissionWindow); ok {
				window.DrawWindowUI(gm, g)
			}
		}
		return
	}

	for _, window := range g.windows {
		window.DrawWindowUI(gm, g)
	}
}

func (g *Game) UpdateWindowSizes(gm *core.GameManager) {
	for _, window := range g.windows {
		window.UpdateWindowSize(gm)
	}
}

func (g *Game) updateCrewSystems() {
	healingPerTick := g.Train.MedbayHealingPerTick()
	damagePerTick := g.Train.LifeSupportDamagePerTick()
	crewCounts := g.stationedCrewCounts()

	g.crewSystemTickTimer += rl.GetFrameTime()
	for g.crewSystemTickTimer >= 1.0 {
		g.crewSystemTickTimer -= 1.0

		for _, character := range g.Train.Characters {
			room, ok := g.Train.GetRoom(character.Pos.RoomId)
			if !ok {
				continue
			}

			if character.IsMoving {
				continue
			}

			if room.System.Type == SystemMedbay && room.IsOperational() {
				character.Health += healingPerTick
				if character.Health > character.MaxHealth {
					character.Health = character.MaxHealth
				}
			}

			if damagePerTick > 0 {
				character.Health -= damagePerTick
				if character.Health < 0 {
					character.Health = 0
				}
			}
		}

		for roomID, crewSupport := range crewCounts {
			effectiveSupport := crewSupport
			if effectiveSupport > 2 {
				effectiveSupport = 2
			}

			if g.resolveThreatInRoom(roomID, effectiveSupport) {
				continue
			}

			room, ok := g.Train.GetRoom(roomID)
			if !ok || room.Health >= room.MaxHealth {
				continue
			}

			room.Health++
			if room.Health > room.MaxHealth {
				room.Health = room.MaxHealth
			}
		}
	}
}

func (g *Game) stationedCrewCounts() map[int]int {
	counts := make(map[int]int)
	for _, character := range g.Train.Characters {
		if character.IsMoving {
			continue
		}

		counts[character.Pos.RoomId]++
	}

	return counts
}

func (g *Game) cartLabel(roomID int) string {
	if roomID < 0 || roomID >= len(g.Train.Rooms) {
		return "Cart ?"
	}

	return fmt.Sprintf("Cart %s", string(g.Train.Rooms[roomID].GetRune()))
}

func (g *Game) CrewSupportSummaryText() string {
	parts := make([]string, 0, len(g.Train.Characters))
	for _, character := range g.Train.Characters {
		if character.IsMoving {
			continue
		}

		parts = append(parts, fmt.Sprintf("%s -> %s", character.Name, g.cartLabel(character.Pos.RoomId)))
	}

	if len(parts) == 0 {
		return "none"
	}

	return strings.Join(parts, "  |  ")
}

func (g *Game) CrewSupportSummaryLines() []string {
	parts := make([]string, 0, len(g.Train.Characters))
	for _, character := range g.Train.Characters {
		if character.IsMoving {
			continue
		}

		parts = append(parts, fmt.Sprintf("%s -> %s", character.Name, g.cartLabel(character.Pos.RoomId)))
	}

	if len(parts) == 0 {
		return []string{"none"}
	}

	return parts
}

func (g *Game) SetCombatStatus(lines ...string) {
	g.combatStatusLines = g.combatStatusLines[:0]
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		g.combatStatusLines = append(g.combatStatusLines, line)
	}

	if len(g.combatStatusLines) == 0 {
		g.combatStatusLines = []string{"No combat actions yet."}
	}
}

func (g *Game) isGameOverModalActive() bool {
	return g.Wave.Failed
}

func (g *Game) ThreatSummaryLines() []string {
	if len(g.ActiveThreats) == 0 {
		return []string{"none"}
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
		return []string{"none"}
	}

	return parts
}

func (g *Game) isMissionBriefingActive() bool {
	return g.MissionBriefingActive
}
