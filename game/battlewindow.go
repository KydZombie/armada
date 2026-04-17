package game

import (
	"fmt"

	"github.com/KydZombie/armada/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type BattleWindow struct {
	core.BaseWindow[Game]

	enemyTexture rl.Texture2D
}

func NewBattleWindow(sizeFunc func(gm *core.GameManager) rl.Rectangle, gm *core.GameManager) *BattleWindow {
	return &BattleWindow{
		BaseWindow:   core.NewBaseWindow[Game](sizeFunc, gm, true),
		enemyTexture: rl.LoadTexture("assets/enemy.png"),
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

	// Keep the text near the top-left corner with a small amount of padding
	// so it is easy to read and does not touch the edges of the window.
	textX := int32(bounds.X) + 12
	textY := int32(bounds.Y) + 12

	// Use padding from the top and right edges so the enemy image sits on
	// the right side of the panel instead of overlapping the text block.
	topPadding := float32(12)
	rightPadding := float32(12)

	// Make the image square and base its size on the panel height so it
	// scales with the battle window. Using a smaller fraction keeps it from
	// becoming too large for short windows.
	imageSize := bounds.Height - topPadding*2
	if imageSize > bounds.Height*0.6 {
		imageSize = bounds.Height * 0.6
	}

	// Position the image box relative to the battle window bounds.
	// This keeps the box anchored to the bottom-right corner with
	// consistent padding from the edges of the battle window.
	holderBounds := rl.Rectangle{
		X:      bounds.X + bounds.Width - rightPadding - imageSize,
		Y:      bounds.Y + bounds.Height - topPadding - imageSize,
		Width:  imageSize,
		Height: imageSize,
	}
	rl.DrawRectangleRec(holderBounds, rl.DarkGray)
	rl.DrawRectangleLinesEx(holderBounds, 2, rl.White)

	// Draw the placeholder enemy image scaled to fit inside the holder.
	rl.DrawTexturePro(
		b.enemyTexture,
		rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  float32(b.enemyTexture.Width),
			Height: float32(b.enemyTexture.Height),
		},
		holderBounds,
		rl.Vector2{},
		0,
		rl.White,
	)

	// If there is no enemy yet, show a simple placeholder message.
	if state.Enemy == nil {
		rl.DrawText("No enemy", textX, textY, 24, rl.White)
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

	rl.DrawText(nameText, textX, textY, 24, rl.White)
	rl.DrawText(healthText, textX, textY+28, 24, rl.White)
	rl.DrawText(statusText, textX, textY+56, 24, statusColor)
}

func (b BattleWindow) DrawWindowUI(gm *core.GameManager, state *Game) {
}
