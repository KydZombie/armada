package game

import (
	"github.com/KydZombie/armada/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Train struct {
	Health, MaxHealth int
}

func NewTrain() *Train {
	return &Train{
		Health:    100,
		MaxHealth: 100,
	}
}

type TrainWindow struct {
	core.BaseWindow[Game]

	train *Train
}

func NewTrainWindow(sizeFunc func(gm *core.GameManager) rl.Rectangle, gm *core.GameManager, train *Train) *TrainWindow {
	return &TrainWindow{
		BaseWindow: core.NewBaseWindow[Game](sizeFunc, gm, true),
		train:      train,
	}
}

func (t TrainWindow) HandleInput(gm *core.GameManager, state *Game) bool {
	return false
}

func (t TrainWindow) UpdateWindow(gm *core.GameManager, state *Game) {

}

func (t TrainWindow) DrawWindow(gm *core.GameManager, state *Game) {
	rl.DrawRectangleRec(t.GetBounds(), rl.Red)
}

func (t TrainWindow) DrawWindowUI(gm *core.GameManager, state *Game) {
}
