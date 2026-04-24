package game

import (
	"fmt"

	"github.com/KydZombie/armada/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type combatantSnapshot struct {
	Name        string
	Health      int
	MaxHealth   int
	AttackPower int
	Alive       bool
}

type BattleWindow struct {
	core.BaseWindow[Game]

	enemyTexture rl.Texture2D
}

func NewBattleWindow(sizeFunc func(gm *core.GameManager) rl.Rectangle, gm *core.GameManager) *BattleWindow {
	return &BattleWindow{
		BaseWindow:   core.NewBaseWindow[Game](sizeFunc, gm, true),
		enemyTexture: rl.LoadTexture("assets/wormy.png"),
	}
}

func (b BattleWindow) HandleInput(gm *core.GameManager, state *Game) bool {
	return false
}

func (b BattleWindow) UpdateWindow(gm *core.GameManager, state *Game) {
}

func (b BattleWindow) drawPlaceholderBar(bounds rl.Rectangle, label string, fillColor rl.Color, textSize int32) {
	rl.DrawRectangleRec(bounds, rl.DarkGray)
	rl.DrawRectangleLinesEx(bounds, 2, rl.White)

	labelText := fmt.Sprintf("%s: pending", label)
	textWidth := rl.MeasureText(labelText, textSize)
	textX := int32(bounds.X + (bounds.Width-float32(textWidth))/2)
	textY := int32(bounds.Y + (bounds.Height-float32(textSize))/2)
	rl.DrawText(labelText, textX, textY, textSize, fillColor)
}

func (b BattleWindow) drawEnemySidebar(sidebarBounds rl.Rectangle, snapshot combatantSnapshot) {
	rl.DrawRectangleRec(sidebarBounds, rl.Fade(rl.Black, 0.35))
	rl.DrawRectangleLinesEx(sidebarBounds, 2, rl.White)

	padding := float32(12)
	textX := int32(sidebarBounds.X + padding)
	textY := int32(sidebarBounds.Y + padding)

	rl.DrawText(snapshot.Name, textX, textY, 22, rl.White)

	statusText := "Status: Active"
	statusColor := rl.White
	if !snapshot.Alive {
		statusText = "Status: Defeated"
		statusColor = rl.Yellow
	}
	rl.DrawText(statusText, textX, textY+24, 18, statusColor)

	imageWidth := sidebarBounds.Width - padding*2
	imageHeight := sidebarBounds.Height * 0.34
	imageBounds := rl.Rectangle{
		X:      sidebarBounds.X + padding,
		Y:      sidebarBounds.Y + 58,
		Width:  imageWidth,
		Height: imageHeight,
	}
	rl.DrawRectangleRec(imageBounds, rl.DarkGray)
	rl.DrawRectangleLinesEx(imageBounds, 2, rl.White)
	rl.DrawTexturePro(
		b.enemyTexture,
		rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  float32(b.enemyTexture.Width),
			Height: float32(b.enemyTexture.Height),
		},
		imageBounds,
		rl.Vector2{},
		0,
		rl.White,
	)

	barX := sidebarBounds.X + padding
	barWidth := sidebarBounds.Width - padding*2
	barHeight := float32(24)
	healthBarY := imageBounds.Y + imageBounds.Height + 14
	drawStatBar(
		rl.Rectangle{
			X:      barX,
			Y:      healthBarY,
			Width:  barWidth,
			Height: barHeight,
		},
		"Health",
		snapshot.Health,
		snapshot.MaxHealth,
		rl.Red,
		16,
	)

	damageBarY := healthBarY + barHeight + 10
	maxDamageDisplay := snapshot.AttackPower
	if maxDamageDisplay < 5 {
		maxDamageDisplay = 5
	}
	drawStatBar(
		rl.Rectangle{
			X:      barX,
			Y:      damageBarY,
			Width:  barWidth,
			Height: barHeight,
		},
		"Damage",
		snapshot.AttackPower,
		maxDamageDisplay,
		rl.Orange,
		16,
	)

	cooldownBarY := damageBarY + barHeight + 10
	b.drawPlaceholderBar(
		rl.Rectangle{
			X:      barX,
			Y:      cooldownBarY,
			Width:  barWidth,
			Height: barHeight,
		},
		"Cooldown",
		rl.SkyBlue,
		16,
	)
}

func (b BattleWindow) DrawWindow(gm *core.GameManager, state *Game) {
	bounds := b.GetBounds()
	rl.DrawRectangleRec(bounds, rl.Red)

	padding := float32(12)
	panelGap := float32(12)
	sidebarWidth := bounds.Width * 0.28
	contentWidth := bounds.Width - sidebarWidth - padding*3

	infoPanel := rl.Rectangle{
		X:      bounds.X + padding,
		Y:      bounds.Y + padding,
		Width:  contentWidth,
		Height: bounds.Height - padding*2,
	}
	sidebarBounds := rl.Rectangle{
		X:      infoPanel.X + infoPanel.Width + panelGap,
		Y:      bounds.Y + padding,
		Width:  sidebarWidth - panelGap,
		Height: bounds.Height - padding*2,
	}

	rl.DrawRectangleRec(infoPanel, rl.Fade(rl.Black, 0.35))
	rl.DrawRectangleLinesEx(infoPanel, 2, rl.White)

	textX := int32(infoPanel.X + padding)
	textY := int32(infoPanel.Y + padding)
	rl.DrawText(
		fmt.Sprintf("Train Hull: %d/%d", state.Train.Health, state.Train.MaxHealth),
		textX,
		textY,
		24,
		rl.White,
	)
	rl.DrawText(
		fmt.Sprintf("Train Damage: %d", state.Train.TotalAttackPower()),
		textX,
		textY+30,
		20,
		rl.White,
	)
	rl.DrawText(
		fmt.Sprintf("Shields: %d   Evasion: %d%%", state.Train.ShieldLayers(), state.Train.EvasionChance()),
		textX,
		textY+56,
		18,
		rl.White,
	)
	rl.DrawText(
		fmt.Sprintf("Medbay Heal: %d/tick", state.Train.MedbayHealingPerTick()),
		textX,
		textY+80,
		18,
		rl.White,
	)
	lifeSupportText := "Life Support: Online"
	lifeSupportColor := rl.White
	if !state.Train.LifeSupportOperational() {
		lifeSupportText = fmt.Sprintf("Life Support: Offline (%d/tick)", state.Train.LifeSupportDamagePerTick())
		lifeSupportColor = rl.Orange
	}
	rl.DrawText(
		lifeSupportText,
		textX,
		textY+104,
		18,
		lifeSupportColor,
	)

	if state.Enemy == nil {
		rl.DrawText("No enemy", textX, textY+136, 24, rl.White)
		return
	}

	enemySnapshot := combatantSnapshot{
		Name:        state.Enemy.Name(),
		Health:      state.Enemy.Health(),
		MaxHealth:   state.Enemy.MaxHealth(),
		AttackPower: state.Enemy.Attack(),
		Alive:       state.Enemy.Alive(),
	}
	b.drawEnemySidebar(sidebarBounds, enemySnapshot)
}

func (b BattleWindow) DrawWindowUI(gm *core.GameManager, state *Game) {
}
