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
}

var (
	battleBgColor        = rl.Black
	battlePanelColor     = rl.Black
	battleInsetColor     = rl.Black
	battleAccentColor    = rl.Gray
	battleBorderColor    = rl.White
	battleEnemyGlowColor = rl.Gray
)

func NewBattleWindow(sizeFunc func(gm *core.GameManager) rl.Rectangle, gm *core.GameManager) *BattleWindow {
	return &BattleWindow{
		BaseWindow: core.NewBaseWindow[Game](sizeFunc, gm, true),
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

func drawFittedText(gm *core.GameManager, text string, x int32, y int32, maxWidth float32, preferredSize int32, minSize int32, color rl.Color) int32 {
	textSize := fitTextSize(text, maxWidth, preferredSize, minSize)
	rl.DrawTextEx(
		gm.Fonts["ec"],
		text,
		rl.NewVector2(float32(x), float32(y)),
		float32(textSize),
		2,
		color,
	)
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

func (b BattleWindow) drawEnemySidebar(gm *core.GameManager, sidebarBounds rl.Rectangle, snapshot combatantSnapshot) {
	padding := float32(12)
	barHeight := float32(24)
	barGap := float32(10)
	headerBounds := rl.Rectangle{
		X:      sidebarBounds.X + padding,
		Y:      sidebarBounds.Y + padding,
		Width:  sidebarBounds.Width - padding*2,
		Height: 56,
	}

	textX := int32(headerBounds.X + 12)
	textY := int32(headerBounds.Y + 8)
	textWidth := headerBounds.Width - 24
	titleSize := fitTextSize(snapshot.Name, textWidth, 26, 14)
	rl.DrawTextEx(
		gm.Fonts["ec"],
		snapshot.Name,
		rl.NewVector2(float32(textX), float32(textY)),
		float32(titleSize),
		2,
		rl.White,
	)

	statusText := "Status: Active"
	statusColor := rl.LightGray
	if !snapshot.Alive {
		statusText = "Status: Defeated"
		statusColor = rl.Yellow
	}
	statusY := textY + titleSize + 4
	statusSize := fitTextSize(statusText, textWidth, 18, 12)
	rl.DrawTextEx(
		gm.Fonts["ec"],
		statusText,
		rl.NewVector2(float32(textX), float32(statusY)),
		float32(statusSize),
		2,
		statusColor,
	)

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

	textureWidth := float32(gm.Textures["enemy"].Width)
	textureHeight := float32(gm.Textures["enemy"].Height)
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

	rl.DrawTexture(gm.Textures["enemyF"], 0, 0, rl.White)
}

func (b BattleWindow) DrawWindow(gm *core.GameManager, state *Game) {
	enemyPanel := rl.Rectangle{
		X:      544,
		Y:      360,
		Width:  192,
		Height: 208,
	}

	if state.Enemy == nil {
		drawPanelCard(enemyPanel, battlePanelColor, rl.Black)
		drawFittedText(gm, "No enemy in range", int32(enemyPanel.X+12), int32(enemyPanel.Y+110), enemyPanel.Width-24, 28, 16, rl.White)
	} else {
		enemySnapshot := combatantSnapshot{
			Name:         state.Enemy.Name(),
			Health:       state.Enemy.Health(),
			MaxHealth:    state.Enemy.MaxHealth(),
			ShieldLayers: state.Enemy.ShieldLayers(),
			AttackPower:  state.Enemy.Attack(),
			Alive:        state.Enemy.Alive(),
		}
		b.drawEnemySidebar(gm, enemyPanel, enemySnapshot)
	}

	rl.DrawTexture(gm.Textures["terminal"], 0, 0, rl.White)
	rl.DrawTexture(gm.Textures["layout"], 0, 0, rl.White)

	rl.DrawCircleGradient(
		0, 0,
		500,
		rl.Color{R: 0, G: 0, B: 0, A: 200},
		rl.Color{R: 0, G: 0, B: 0, A: 0},
	)
	rl.DrawCircleGradient(
		gm.NativeWidth, 0,
		500,
		rl.Color{R: 0, G: 0, B: 0, A: 200},
		rl.Color{R: 0, G: 0, B: 0, A: 0},
	)
	rl.DrawCircleGradient(
		0, gm.NativeHeight,
		500,
		rl.Color{R: 0, G: 0, B: 0, A: 200},
		rl.Color{R: 0, G: 0, B: 0, A: 0},
	)
	rl.DrawCircleGradient(
		gm.NativeWidth, gm.NativeHeight,
		500,
		rl.Color{R: 0, G: 0, B: 0, A: 200},
		rl.Color{R: 0, G: 0, B: 0, A: 0},
	)
}

func (b BattleWindow) DrawWindowUI(gm *core.GameManager, state *Game) {
}
