package game

import (
	"github.com/KydZombie/armada/core"
	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	moving bool
	pos    rl.Vector2

	windows []core.Window[Game]
}

func NewGameScreen(gm *core.GameManager) *Game {
	gs := &Game{
		moving: false,
		pos:    rl.Vector2{X: 200, Y: 50},

		windows: []core.Window[Game]{},
	}

	gs.windows = append(gs.windows, NewTerminalWindow(
		rl.Rectangle{
			X:      float32(gm.ScreenWidth) / 2.0,
			Y:      float32(gm.ScreenHeight) / 2.0,
			Width:  float32(gm.ScreenWidth) / 2.0,
			Height: float32(gm.ScreenHeight) / 2.0,
		},
		initializeCommands(),
	))

	return gs
}

func (g *Game) UpdateScreen(gm *core.GameManager) {
	if g.moving {
		g.pos.X += 100.0 * gm.DeltaTime
	}

	inputCaptured := false

	for _, window := range g.windows {
		// If a window captures the input, other windows should not read any input
		if !inputCaptured && window.HandleInput(gm, g) {
			inputCaptured = true
		}

		window.UpdateWindow(gm, g)
	}
}

func (g *Game) DrawScreen(gm *core.GameManager) {
	rl.ClearBackground(rl.DarkBlue)
	rl.DrawRectangleV(g.pos, rl.Vector2{X: 50, Y: 50}, rl.Red)

	for _, window := range g.windows {
		window.DrawWindow(gm, g)
	}
}

func (g *Game) DrawScreenUI(gm *core.GameManager) {
	var buttonText string
	if g.moving {
		buttonText = "Stop moving"
	} else {
		buttonText = "Start moving"
	}

	if rg.Button(rl.Rectangle{
		X:      0,
		Y:      50,
		Width:  100,
		Height: 60,
	}, buttonText) {
		g.moving = !g.moving
	}

	for _, window := range g.windows {
		window.DrawWindow(gm, g)
	}
}
