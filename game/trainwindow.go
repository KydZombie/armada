package game

import (
	"github.com/KydZombie/armada/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type TrainWindow struct {
	core.BaseWindow[Game]
}

func NewTrainWindow(sizeFunc func(gm *core.GameManager) rl.Rectangle, gm *core.GameManager) *TrainWindow {
	return &TrainWindow{
		BaseWindow: core.NewBaseWindow[Game](sizeFunc, gm, true),
	}
}

func (t TrainWindow) HandleInput(gm *core.GameManager, state *Game) bool {
	return false
}

func (t TrainWindow) UpdateWindow(gm *core.GameManager, state *Game) {

}

func (t TrainWindow) DrawWindow(gm *core.GameManager, state *Game) {
	rl.DrawRectangleRec(t.GetBounds(), rl.Blue)
}

func (t TrainWindow) DrawWindowUI(gm *core.GameManager, state *Game) {
}
