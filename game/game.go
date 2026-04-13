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

	const windowMargin = 16.0

	terminal := NewTerminalWindow(
		func(gm *core.GameManager) rl.Rectangle {
			return rl.Rectangle{
				X:      float32(gm.ScreenWidth)/2.0 + windowMargin,
				Y:      float32(gm.ScreenHeight)/2.0 + windowMargin,
				Width:  float32(gm.ScreenWidth)/2.0 - windowMargin*2,
				Height: float32(gm.ScreenHeight)/2.0 - windowMargin,
			}
		},
		gm,
		initializeCommands(),
	)

	gs.windows = append(gs.windows, terminal)

	gs.windows = append(gs.windows, NewTrainWindow(
		func(gm *core.GameManager) rl.Rectangle {
			return rl.Rectangle{
				X:      windowMargin,
				Y:      windowMargin,
				Width:  float32(gm.ScreenWidth)/2.0 - windowMargin,
				Height: float32(gm.ScreenHeight) - windowMargin*2,
			}
		},
		gm,
	))

	gs.windows = append(gs.windows, NewBattleWindow(
		func(gm *core.GameManager) rl.Rectangle {
			if terminal.IsVisible() {
				return rl.Rectangle{
					X:      float32(gm.ScreenWidth)/2.0 + windowMargin,
					Y:      windowMargin,
					Width:  float32(gm.ScreenWidth)/2.0 - windowMargin*2,
					Height: float32(gm.ScreenHeight)/2.0 - windowMargin,
				}
			} else {
				return rl.Rectangle{
					X:      float32(gm.ScreenWidth)/2.0 + windowMargin,
					Y:      windowMargin,
					Width:  float32(gm.ScreenWidth)/2.0 - windowMargin*2,
					Height: float32(gm.ScreenHeight) - windowMargin*2,
				}
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

func (g *Game) UpdateWindowSizes(gm *core.GameManager) {
	for _, window := range g.windows {
		window.UpdateWindowSize(gm)
	}
}
