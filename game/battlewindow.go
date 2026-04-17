package game

import (
	"fmt"

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
	// Draw the battle window background first so the enemy text sits
	// on a readable panel.
	bounds := b.GetBounds()
	rl.DrawRectangleRec(bounds, rl.Red)

	// If there is no enemy yet, show a simple placeholder message.
	if state.Enemy == nil {
		rl.DrawText("No enemy", int32(bounds.X)+8, int32(bounds.Y)+8, 24, rl.White)
		return
	}

	// Show the enemy's basic information so the prototype has visible
	// feedback when the terminal attack command is used.
	nameText := fmt.Sprint("Enemy: ", state.Enemy.Name())
	healthText := fmt.Sprintf("Health: %d/%d", state.Enemy.Health(), state.Enemy.MaxHealth())

	statusText := "Status: Alive"
	statusColor := rl.White
	if !state.Enemy.Alive() {
		statusText = "Status: Defeated"
		statusColor = rl.Yellow
	}

	textX := int32(bounds.X) + 8
	textY := int32(bounds.Y) + 8

	rl.DrawText(nameText, textX, textY, 24, rl.White)
	rl.DrawText(healthText, textX, textY+28, 24, rl.White)
	rl.DrawText(statusText, textX, textY+56, 24, statusColor)
}

func (b BattleWindow) DrawWindowUI(gm *core.GameManager, state *Game) {
}
