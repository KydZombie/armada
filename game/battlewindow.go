package game

import (
	"fmt"

	"github.com/KydZombie/armada/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type BattleWindow struct {
	core.BaseWindow[Game]

	enemyTextures map[rune]rl.Texture2D
}

func NewBattleWindow(sizeFunc func(gm *core.GameManager) rl.Rectangle, gm *core.GameManager) *BattleWindow {
	enemyTextures := map[rune]rl.Texture2D{
		'A': rl.LoadTexture("assets/centipede.png"),
		'B': rl.LoadTexture("assets/spidy.png"),
		'C': rl.LoadTexture("assets/wormy.png"),
	}

	for roomLabel, texture := range enemyTextures {
		if texture.ID == 0 {
			gm.ErrLog.Printf("battle enemy texture failed to load for room %c", roomLabel)
			continue
		}

		gm.Log.Printf(
			"battle enemy texture loaded for room %c: %dx%d",
			roomLabel,
			texture.Width,
			texture.Height,
		)
	}

	return &BattleWindow{
		BaseWindow: core.NewBaseWindow[Game](sizeFunc, gm, true),
		enemyTextures: enemyTextures,
	}
}

func (b BattleWindow) HandleInput(gm *core.GameManager, state *Game) bool {
	return false
}

func (b BattleWindow) UpdateWindow(gm *core.GameManager, state *Game) {
}

func (b BattleWindow) DrawWindow(gm *core.GameManager, state *Game) {
	bounds := b.GetBounds()
	rl.DrawRectangleRec(bounds, rl.Red)

	topPadding := float32(12)
	rightPadding := float32(12)
	popupWidth := float32(220)
	popupHeight := float32(160)
	maxPopupWidth := bounds.Width - rightPadding*2
	maxPopupHeight := bounds.Height - topPadding*2

	if popupWidth > maxPopupWidth {
		popupWidth = maxPopupWidth
	}
	if popupHeight > maxPopupHeight {
		popupHeight = maxPopupHeight
	}
	if popupWidth <= 0 || popupHeight <= 0 {
		popupWidth = 0
		popupHeight = 0
	}

	selectedLabelBounds := rl.Rectangle{
		X:      bounds.X + bounds.Width - rightPadding - popupWidth,
		Y:      bounds.Y + bounds.Height - topPadding - popupHeight,
		Width:  popupWidth,
		Height: popupHeight,
	}
	selectedTargetLabel := rune('A' + state.SelectedRoom)
	selectedTargetText := fmt.Sprintf("Selected: [%c]", selectedTargetLabel)
	selectedTargetDisplayText := selectedTargetText
	if state.SelectionPopupFrames > 0 && state.SelectionPopupText != "" {
		selectedTargetDisplayText = state.SelectionPopupText
	}
	var selectedEnemy Enemy
	if state.SelectedRoom >= 0 && state.SelectedRoom < len(state.RoomEnemies) {
		selectedEnemy = state.RoomEnemies[state.SelectedRoom]
	}
	displayEnemy := selectedEnemy
	if displayEnemy == nil {
		displayEnemy = state.Enemy
	}
	cooldownBarBounds := rl.Rectangle{}
	cooldownFillBounds := rl.Rectangle{}
	drawCooldownBar := false
	textX := int32(bounds.X) + 12
	textY := int32(bounds.Y) + 12

	if displayEnemy == nil {
		rl.DrawText("No enemy", textX, textY, 24, rl.White)
	} else {
		nameText := fmt.Sprint("Enemy: ", displayEnemy.Name())
		healthText := fmt.Sprintf("Health: %d/%d", displayEnemy.Health(), displayEnemy.MaxHealth())

		statusText := "Status: Alive"
		statusColor := rl.White
		if !displayEnemy.Alive() {
			statusText = "Status: Defeated"
			statusColor = rl.Yellow
		}
		statusTextY := textY + 56

		rl.DrawText(nameText, textX, textY, 24, rl.White)
		rl.DrawText(healthText, textX, textY+28, 24, rl.White)
		rl.DrawText(statusText, textX, statusTextY, 24, statusColor)

		enemy, ok := displayEnemy.(*BasicEnemy)
		if ok {
			if enemy.attackCooldown > 0 {
				timerFillRatio := 1 - (float32(enemy.attackTimer) / float32(enemy.attackCooldown))
				if timerFillRatio < 0 {
					timerFillRatio = 0
				}
				if timerFillRatio > 1 {
					timerFillRatio = 1
				}

				timerBarHeight := float32(8)
				timerBarGap := float32(12)
				timerBarWidth := popupWidth
				timerBarX := selectedLabelBounds.X
				timerBarY := selectedLabelBounds.Y - timerBarGap - timerBarHeight

				if timerBarWidth <= 0 {
					timerBarWidth = 100
					timerBarX = bounds.X + bounds.Width - rightPadding - timerBarWidth
				}
				if timerBarY < bounds.Y+topPadding {
					timerBarY = bounds.Y + topPadding
				}

				cooldownBarBounds = rl.Rectangle{
					X:      timerBarX,
					Y:      timerBarY,
					Width:  timerBarWidth,
					Height: timerBarHeight,
				}
				cooldownFillBounds = rl.Rectangle{
					X:      cooldownBarBounds.X,
					Y:      cooldownBarBounds.Y,
					Width:  cooldownBarBounds.Width * timerFillRatio,
					Height: cooldownBarBounds.Height,
				}
				drawCooldownBar = true
			}

			partTextX := textX
			partStartY := textY + 112
			barWidth := float32(100)
			barHeight := float32(10)
			partSpacing := int32(34)
			maxPartWidth := int32(bounds.Width) - 24
			if popupWidth > 0 {
				maxPartWidth = int32(selectedLabelBounds.X) - partTextX - 12
			}

			for i, part := range enemy.Parts {
				if part == nil {
					continue
				}

				partY := partStartY + int32(i)*partSpacing
				barY := partY + 18

				if float32(barY)+barHeight > bounds.Y+bounds.Height-12 {
					break
				}

				rl.DrawText(part.Name, partTextX, partY, 20, rl.White)

				currentBarWidth := barWidth
				if maxPartWidth > 0 && currentBarWidth > float32(maxPartWidth) {
					currentBarWidth = float32(maxPartWidth)
				}

				barBounds := rl.Rectangle{
					X:      float32(partTextX),
					Y:      float32(barY),
					Width:  currentBarWidth,
					Height: barHeight,
				}
				rl.DrawRectangleRec(barBounds, rl.DarkGray)

				fillRatio := float32(0)
				if part.MaxHealth > 0 {
					fillRatio = float32(part.Health) / float32(part.MaxHealth)
				}
				if fillRatio < 0 {
					fillRatio = 0
				}
				if fillRatio > 1 {
					fillRatio = 1
				}

				fillBounds := rl.Rectangle{
					X:      barBounds.X,
					Y:      barBounds.Y,
					Width:  barBounds.Width * fillRatio,
					Height: barBounds.Height,
				}
				rl.DrawRectangleRec(fillBounds, rl.Green)
			}
		}
	}

	if drawCooldownBar {
		rl.DrawText(
			"Attack Cooldown",
			int32(cooldownBarBounds.X),
			int32(cooldownBarBounds.Y)-24,
			20,
			rl.White,
		)
		rl.DrawRectangleRec(cooldownBarBounds, rl.DarkGray)
		rl.DrawRectangleRec(cooldownFillBounds, rl.Orange)
	}

	if popupWidth > 0 && popupHeight > 0 {
		rl.DrawRectangleRec(selectedLabelBounds, rl.DarkGray)
		rl.DrawRectangleLinesEx(selectedLabelBounds, 2, rl.White)

		rl.DrawText(
			selectedTargetDisplayText,
			int32(selectedLabelBounds.X)+8,
			int32(selectedLabelBounds.Y)+10,
			20,
			rl.White,
		)

		b.drawPopupTargetImage(selectedLabelBounds, selectedTargetLabel)
	}
}

func (b BattleWindow) DrawWindowUI(gm *core.GameManager, state *Game) {
}

func (b BattleWindow) drawPopupTargetImage(popupBounds rl.Rectangle, roomLabel rune) {
	texture, ok := b.enemyTextures[roomLabel]
	if !ok || texture.ID == 0 {
		return
	}

	const imagePadding float32 = 8.0
	const textHeight float32 = 24.0

	destBounds := rl.Rectangle{
		X:      popupBounds.X + imagePadding,
		Y:      popupBounds.Y + 10 + textHeight + imagePadding,
		Width:  popupBounds.Width - imagePadding*2,
		Height: popupBounds.Height - textHeight - imagePadding*3 - 10,
	}
	if destBounds.Width <= 0 || destBounds.Height <= 0 {
		return
	}

	rl.DrawTexturePro(
		texture,
		rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  float32(texture.Width),
			Height: float32(texture.Height),
		},
		destBounds,
		rl.Vector2{},
		0,
		rl.White,
	)
}
