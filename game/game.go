package game

import (
	"github.com/KydZombie/armada/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	Train *Train

	windows []core.Window[Game]
}

func NewGameScreen(gm *core.GameManager) *Game {
	train := NewTrain()

	gs := &Game{
		Train: train,

		windows: []core.Window[Game]{},
	}

	gs.windows = append(gs.windows, NewTerminalWindow(
		func(gm *core.GameManager) rl.Rectangle {
			return rl.Rectangle{
				X:      float32(gm.ScreenWidth) / 2.0,
				Y:      float32(gm.ScreenHeight) / 2.0,
				Width:  float32(gm.ScreenWidth) / 2.0,
				Height: float32(gm.ScreenHeight) / 2.0,
			}
		},
		gm,
		initializeCommands(),
	))

	gs.windows = append(gs.windows, NewTrainWindow(
		func(gm *core.GameManager) rl.Rectangle {
			return rl.Rectangle{
				X:      32.0,
				Y:      32.0,
				Width:  float32(gm.ScreenWidth)/2.0 - 32.0,
				Height: float32(gm.ScreenHeight) - 64.0,
			}
		},
		gm,
		train,
	))

	return gs
}

func (g *Game) ResizeScreen(gm *core.GameManager) {
	for _, w := range g.windows {
		w.ResizeWindow(gm)
	}
}

func (g *Game) UpdateScreen(gm *core.GameManager) {
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

	for _, window := range g.windows {
		window.DrawWindow(gm, g)
	}
}

func (g *Game) DrawScreenUI(gm *core.GameManager) {
	//var buttonText string
	//if g.moving {
	//	buttonText = "Stop moving"
	//} else {
	//	buttonText = "Start moving"
	//}
	//
	//if rg.Button(rl.Rectangle{
	//	X:      0,
	//	Y:      50,
	//	Width:  100,
	//	Height: 60,
	//}, buttonText) {
	//	g.moving = !g.moving
	//}

	for _, window := range g.windows {
		window.DrawWindow(gm, g)
	}
}
