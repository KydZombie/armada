package game

import (
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

	SelectedWeaponIndex int
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
	}

	const windowMargin = 16.0

	terminal := NewTerminalWindow(
		func(gm *core.GameManager) rl.Rectangle {
			return rl.Rectangle{
				X:      windowMargin,
				Y:      float32(gm.ScreenHeight)/2.0 + windowMargin,
				Width:  float32(gm.ScreenWidth)/2.0 - windowMargin,
				Height: float32(gm.ScreenHeight)/2.0 - windowMargin*2,
			}
		},
		gm,
		gs.Terminal,
	)

	gs.windows = append(gs.windows, terminal)

	gs.windows = append(gs.windows, NewTrainWindow(
		func(gm *core.GameManager) rl.Rectangle {
			return rl.Rectangle{
				X:      windowMargin,
				Y:      windowMargin,
				Width:  float32(gm.ScreenWidth) - windowMargin*2,
				Height: float32(gm.ScreenHeight)/2.0 - windowMargin,
			}
		},
		gm,
	))

	gs.windows = append(gs.windows, NewBattleWindow(
		func(gm *core.GameManager) rl.Rectangle {
			return rl.Rectangle{
				X:      float32(gm.ScreenWidth)/2.0 + windowMargin,
				Y:      float32(gm.ScreenHeight)/2.0 + windowMargin,
				Width:  float32(gm.ScreenWidth)/2.0 - windowMargin*2,
				Height: float32(gm.ScreenHeight)/2.0 - windowMargin*2,
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

	if g.Train.WeaponsOperational() {
		g.Train.AdvanceWeaponCooldowns(deltaSeconds)
	}
	g.Train.UpdateCombatState(deltaSeconds)
	if g.Enemy != nil {
		g.Enemy.UpdateCombatState(deltaSeconds)
	}
	g.updateCrewSystems()
}

func (g *Game) DrawScreen(gm *core.GameManager) {
	rl.ClearBackground(rl.DarkBlue)

	for _, window := range g.windows {
		window.DrawWindow(gm, g)
	}
}

func (g *Game) DrawScreenUI(gm *core.GameManager) {
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

	g.crewSystemTickTimer += rl.GetFrameTime()
	if g.crewSystemTickTimer < 1.0 {
		return
	}

	g.crewSystemTickTimer = 0
	for _, character := range g.Train.Characters {
		room, ok := g.Train.GetRoom(character.Pos.RoomId)
		if !ok {
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
