package game

import (
	"github.com/KydZombie/armada/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type BattleWindow struct {
	core.BaseWindow[Game]
}

func NewBattleWindow(sizeFunc func(gm *core.GameManager) rl.Rectangle, gm *core.GameManager) *BattleWindow {
	return &BattleWindow{
		BaseWindow: core.NewBaseWindow[Game](sizeFunc, gm, true),
	}
}

func (b BattleWindow) HandleInput(gm *core.GameManager, state *Game) bool {
	return false
}

func (b BattleWindow) UpdateWindow(gm *core.GameManager, state *Game) {
}

func (b BattleWindow) DrawWindow(gm *core.GameManager, state *Game) {
	rl.DrawRectangleRec(b.GetBounds(), rl.Red)
}

func (b BattleWindow) DrawWindowUI(gm *core.GameManager, state *Game) {
}
