package game

import (
	"github.com/KydZombie/armada/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type combatantSnapshot struct {
	Name         string
	Health       int
	MaxHealth    int
	ShieldLayers int
	AttackPower  int
	Alive        bool
}

type BattleWindow struct {
	core.BaseWindow[Game]

	enemyTexture rl.Texture2D
}

var (
	battleBgColor        = rl.Red
	battlePanelColor     = rl.NewColor(20, 28, 42, 245)
	battleInsetColor     = rl.NewColor(28, 38, 56, 255)
	battleAccentColor    = rl.NewColor(94, 176, 239, 255)
	battleBorderColor    = rl.NewColor(188, 212, 230, 235)
	battleEnemyGlowColor = rl.NewColor(167, 79, 48, 255)
)

func NewBattleWindow(sizeFunc func(gm *core.GameManager) rl.Rectangle, gm *core.GameManager) *BattleWindow {
	return &BattleWindow{
		BaseWindow:   core.NewBaseWindow[Game](sizeFunc, gm, true),
		enemyTexture: rl.LoadTexture("assets/wormy.png"),
	}
}

func (b BattleWindow) HandleInput(gm *core.GameManager, state *Game) bool {
	if state.isGameOverModalActive() || state.isMissionBriefingActive() {
		return false
	}

	return false
}

func (b BattleWindow) UpdateWindow(gm *core.GameManager, state *Game) {
}

func fitTextSize(text string, maxWidth float32, preferredSize int32, minSize int32) int32 {
	if preferredSize < minSize {
		preferredSize = minSize
	}

	for size := preferredSize; size > minSize; size-- {
		if float32(rl.MeasureText(text, size)) <= maxWidth {
			return size
		}
	}

	return minSize
}

func drawFittedText(text string, x int32, y int32, maxWidth float32, preferredSize int32, minSize int32, color rl.Color) int32 {
	textSize := fitTextSize(text, maxWidth, preferredSize, minSize)
	rl.DrawText(text, x, y, textSize, color)
	return textSize
}

func drawPanelCard(bounds rl.Rectangle, fill rl.Color, border rl.Color) {
	rl.DrawRectangleRec(bounds, fill)
	rl.DrawRectangleLinesEx(bounds, 2, border)
}

func drawSectionHeader(bounds rl.Rectangle, title string, accent rl.Color) {
	drawPanelCard(bounds, battleInsetColor, rl.Fade(accent, 0.8))
	rl.DrawRectangleRec(rl.Rectangle{
		X:      bounds.X,
		Y:      bounds.Y,
		Width:  4,
		Height: bounds.Height,
	}, accent)
	rl.DrawText(title, int32(bounds.X+12), int32(bounds.Y+8), fitTextSize(title, bounds.Width-24, 18, 14), rl.White)
}

func (b BattleWindow) drawEnemySidebar(sidebarBounds rl.Rectangle, snapshot combatantSnapshot) {
	drawPanelCard(sidebarBounds, battlePanelColor, rl.Fade(battleBorderColor, 0.85))

	padding := float32(12)
	barHeight := float32(24)
	barGap := float32(10)
	headerBounds := rl.Rectangle{
		X:      sidebarBounds.X + padding,
		Y:      sidebarBounds.Y + padding,
		Width:  sidebarBounds.Width - padding*2,
		Height: 56,
	}
	drawPanelCard(headerBounds, battleInsetColor, rl.Fade(battleEnemyGlowColor, 0.9))
	rl.DrawRectangleRec(rl.Rectangle{
		X:      headerBounds.X,
		Y:      headerBounds.Y,
		Width:  5,
		Height: headerBounds.Height,
	}, battleEnemyGlowColor)

	textX := int32(headerBounds.X + 12)
	textY := int32(headerBounds.Y + 8)
	textWidth := headerBounds.Width - 24
	titleSize := fitTextSize(snapshot.Name, textWidth, 26, 14)
	rl.DrawText(snapshot.Name, textX, textY, titleSize, rl.White)

	statusText := "Status: Active"
	statusColor := rl.White
	if !snapshot.Alive {
		statusText = "Status: Defeated"
		statusColor = rl.Yellow
	}
	statusY := textY + titleSize + 4
	statusSize := fitTextSize(statusText, textWidth, 18, 12)
	rl.DrawText(statusText, textX, statusY, statusSize, statusColor)

	imageWidth := sidebarBounds.Width - padding*2
	statsHeight := barHeight*3 + barGap*2 + 14
	imageTopY := headerBounds.Y + headerBounds.Height + 10
	imageBottomLimit := sidebarBounds.Y + sidebarBounds.Height - padding - statsHeight
	imageHeight := imageBottomLimit - imageTopY
	if imageHeight < 80 {
		imageHeight = 80
	}
	imageBounds := rl.Rectangle{
		X:      sidebarBounds.X + padding,
		Y:      imageTopY,
		Width:  imageWidth,
		Height: imageHeight,
	}
	drawPanelCard(imageBounds, rl.NewColor(34, 28, 30, 255), rl.Fade(battleEnemyGlowColor, 0.9))

	textureWidth := float32(b.enemyTexture.Width)
	textureHeight := float32(b.enemyTexture.Height)
	drawBounds := imageBounds
	if textureWidth > 0 && textureHeight > 0 {
		textureAspect := textureWidth / textureHeight
		boundsAspect := imageBounds.Width / imageBounds.Height
		if textureAspect > boundsAspect {
			drawBounds.Height = imageBounds.Width / textureAspect
			drawBounds.Y = imageBounds.Y + (imageBounds.Height-drawBounds.Height)/2
		} else {
			drawBounds.Width = imageBounds.Height * textureAspect
			drawBounds.X = imageBounds.X + (imageBounds.Width-drawBounds.Width)/2
		}
	}

	rl.DrawTexturePro(
		b.enemyTexture,
		rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  float32(b.enemyTexture.Width),
			Height: float32(b.enemyTexture.Height),
		},
		drawBounds,
		rl.Vector2{},
		0,
		rl.White,
	)

	barX := sidebarBounds.X + padding
	barWidth := sidebarBounds.Width - padding*2
	statsBounds := rl.Rectangle{
		X:      barX,
		Y:      imageBounds.Y + imageBounds.Height + 10,
		Width:  barWidth,
		Height: sidebarBounds.Y + sidebarBounds.Height - padding - (imageBounds.Y + imageBounds.Height + 10),
	}
	drawPanelCard(statsBounds, battleInsetColor, rl.Fade(battleEnemyGlowColor, 0.6))

	healthBarY := statsBounds.Y + 10
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
		18,
	)

	damageBarY := healthBarY + barHeight + barGap
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
		"Attack",
		snapshot.AttackPower,
		maxDamageDisplay,
		rl.Orange,
		18,
	)

	cooldownBarY := damageBarY + barHeight + barGap
	drawStatBar(
		rl.Rectangle{
			X:      barX,
			Y:      cooldownBarY,
			Width:  barWidth,
			Height: barHeight,
		},
		"Shields",
		snapshot.ShieldLayers,
		4,
		rl.SkyBlue,
		18,
	)
}

func (b BattleWindow) DrawWindow(gm *core.GameManager, state *Game) {
	bounds := b.GetBounds()
	rl.DrawRectangleRec(bounds, battleBgColor)
	rl.DrawRectangleLinesEx(bounds, 2, rl.Fade(battleBorderColor, 0.6))

	padding := float32(12)
	enemyPanel := rl.Rectangle{
		X:      bounds.X + padding,
		Y:      bounds.Y + padding,
		Width:  bounds.Width - padding*2,
		Height: bounds.Height - padding*2,
	}

	if state.Enemy == nil {
		drawPanelCard(enemyPanel, battlePanelColor, rl.Fade(battleBorderColor, 0.85))
		drawFittedText("No enemy in range", int32(enemyPanel.X+12), int32(enemyPanel.Y+16), enemyPanel.Width-24, 28, 16, rl.White)
		drawFittedText("Complete mission goals to spawn the next hostile.", int32(enemyPanel.X+12), int32(enemyPanel.Y+50), enemyPanel.Width-24, 16, 11, rl.LightGray)
		return
	}

	enemySnapshot := combatantSnapshot{
		Name:         state.Enemy.Name(),
		Health:       state.Enemy.Health(),
		MaxHealth:    state.Enemy.MaxHealth(),
		ShieldLayers: state.Enemy.ShieldLayers(),
		AttackPower:  state.Enemy.Attack(),
		Alive:        state.Enemy.Alive(),
	}
	b.drawEnemySidebar(enemyPanel, enemySnapshot)
}

func (b BattleWindow) DrawWindowUI(gm *core.GameManager, state *Game) {
}
