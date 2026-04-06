package game

import (
	"github.com/KydZombie/armada/core"
	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

//goland:noinspection GoNameStartsWithPackageName
type GameScreen struct {
	moving bool
	pos    rl.Vector2
}

func NewGameScreen() *GameScreen {
	return &GameScreen{
		moving: false,
		pos:    rl.Vector2{X: 200, Y: 50},
	}
}

func (g *GameScreen) Update(gm *core.GameManager) {
	if g.moving {
		g.pos.X += 100.0 * gm.DeltaTime
	}
}

func (g *GameScreen) Draw(gm *core.GameManager) {
	rl.ClearBackground(rl.DarkBlue)
	rl.DrawRectangleV(g.pos, rl.Vector2{X: 50, Y: 50}, rl.Red)
}

func (g *GameScreen) DrawUI(gm *core.GameManager) {
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
}
