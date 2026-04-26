package game

import (
	"fmt"
	"strings"

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

func (b BattleWindow) drawWeaponStatusPanel(bounds rl.Rectangle, state *Game) {
	drawPanelCard(bounds, battlePanelColor, rl.Fade(battleBorderColor, 0.85))

	padding := float32(10)
	headerBounds := rl.Rectangle{
		X:      bounds.X + padding,
		Y:      bounds.Y + padding,
		Width:  bounds.Width - padding*2,
		Height: 32,
	}
	drawSectionHeader(headerBounds, "Weapons", rl.Gold)

	barY := headerBounds.Y + headerBounds.Height + 10
	barHeight := float32(20)
	for i, weapon := range state.Train.Weapons {
		rowY := barY + float32(i)*34
		statusColor := rl.Gold
		currentValue := int(weapon.ChargeProgress() * 100)
		maxValue := 100
		label := weapon.Name
		if weapon.Ready() {
			currentValue = 100
			statusColor = rl.Green
			label = fmt.Sprintf("%s ready", weapon.Name)
		} else {
			label = fmt.Sprintf("%s charging %ds", weapon.Name, weapon.CooldownDisplaySeconds())
		}

		drawStatBar(
			rl.Rectangle{
				X:      bounds.X + padding,
				Y:      rowY,
				Width:  bounds.Width - padding*2,
				Height: barHeight,
			},
			label,
			currentValue,
			maxValue,
			statusColor,
			14,
		)
	}
}

func (b BattleWindow) drawCombatStatusPanel(bounds rl.Rectangle, state *Game) {
	drawPanelCard(bounds, battlePanelColor, rl.Fade(battleBorderColor, 0.85))

	padding := float32(10)
	headerBounds := rl.Rectangle{
		X:      bounds.X + padding,
		Y:      bounds.Y + padding,
		Width:  bounds.Width - padding*2,
		Height: 32,
	}
	drawSectionHeader(headerBounds, "Combat Feed", battleAccentColor)

	contentBounds := rl.Rectangle{
		X:      bounds.X + padding,
		Y:      headerBounds.Y + headerBounds.Height + 10,
		Width:  bounds.Width - padding*2,
		Height: bounds.Height - (headerBounds.Height + padding*2 + 10),
	}
	drawPanelCard(contentBounds, rl.Fade(battleInsetColor, 0.92), rl.Fade(battleBorderColor, 0.45))

	textX := int32(contentBounds.X + 12)
	lineY := int32(contentBounds.Y + 12)
	maxTextWidth := contentBounds.Width - 24
	maxTextY := int32(contentBounds.Y + contentBounds.Height - 12)
	for i, line := range state.combatStatusLines {
		if i >= 4 {
			break
		}
		color := rl.LightGray
		bulletColor := rl.LightGray
		if strings.Contains(line, "evades") {
			color = rl.SkyBlue
			bulletColor = rl.SkyBlue
		} else if strings.Contains(line, "shields absorb") {
			color = rl.Gold
			bulletColor = rl.Gold
		} else if strings.Contains(line, "shields and deals no hull damage") {
			color = rl.Gold
			bulletColor = rl.Gold
		} else if strings.Contains(line, "attacks back") {
			color = rl.Orange
			bulletColor = rl.Orange
		}

		lineSize := fitTextSize(line, maxTextWidth, 16, 11)
		wrappedLines := wrapTerminalLine(line, maxTextWidth-18, lineSize)
		for j, wrappedLine := range wrappedLines {
			if lineY+lineSize > maxTextY {
				return
			}
			if j == 0 {
				rl.DrawCircle(int32(contentBounds.X+20), lineY+lineSize/2, 4, bulletColor)
			}
			rl.DrawText(wrappedLine, textX+18, lineY, lineSize, color)
			lineY += lineSize + 2
		}
		lineY += 6
	}
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
	titleSize := fitTextSize(snapshot.Name, textWidth, 22, 12)
	rl.DrawText(snapshot.Name, textX, textY, titleSize, rl.White)

	statusText := "Status: Active"
	statusColor := rl.White
	if !snapshot.Alive {
		statusText = "Status: Defeated"
		statusColor = rl.Yellow
	}
	statusY := textY + titleSize + 4
	statusSize := fitTextSize(statusText, textWidth, 16, 11)
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
		16,
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
		16,
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
		16,
	)
}

func (b BattleWindow) DrawWindow(gm *core.GameManager, state *Game) {
	bounds := b.GetBounds()
	rl.DrawRectangleRec(bounds, battleBgColor)
	rl.DrawRectangleLinesEx(bounds, 2, rl.Fade(battleBorderColor, 0.6))

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

	drawPanelCard(infoPanel, battlePanelColor, rl.Fade(battleBorderColor, 0.85))

	textX := int32(infoPanel.X + padding)
	infoTextWidth := infoPanel.Width - padding*2
	headerBounds := rl.Rectangle{
		X:      infoPanel.X + padding,
		Y:      infoPanel.Y + padding,
		Width:  infoPanel.Width - padding*2,
		Height: 78,
	}
	drawPanelCard(headerBounds, battleInsetColor, rl.Fade(battleAccentColor, 0.85))
	rl.DrawText("Train Status", int32(headerBounds.X+12), int32(headerBounds.Y+8), fitTextSize("Train Status", headerBounds.Width-24, 20, 14), rl.White)

	hullText := fmt.Sprintf("Hull %d/%d", state.Train.Health, state.Train.MaxHealth)
	hullTextSize := drawFittedText(hullText, textX, int32(headerBounds.Y+32), infoTextWidth, 24, 14, rl.White)
	readyDamageText := fmt.Sprintf("Ready Damage %d", state.Train.TotalAttackPower())
	drawFittedText(readyDamageText, textX+220, int32(headerBounds.Y+32), infoTextWidth-220, 18, 12, rl.Gold)
	defenseText := fmt.Sprintf("Shields %d   Evasion %d%%   Weapons %d/%d", state.Train.ShieldLayers(), state.Train.EvasionChance(), state.Train.ReadyWeapons(), len(state.Train.Weapons))
	defenseY := int32(headerBounds.Y+32) + hullTextSize + 4
	drawFittedText(defenseText, textX, defenseY, infoTextWidth, 17, 11, rl.Fade(rl.White, 0.9))

	lifeSupportText := "Life Support: Online"
	lifeSupportColor := rl.White
	if !state.Train.LifeSupportOperational() {
		lifeSupportText = fmt.Sprintf("Life Support: Offline (%d/tick)", state.Train.LifeSupportDamagePerTick())
		lifeSupportColor = rl.Orange
	}

	auxBounds := rl.Rectangle{
		X:      infoPanel.X + padding,
		Y:      headerBounds.Y + headerBounds.Height + 10,
		Width:  infoPanel.Width - padding*2,
		Height: 54,
	}
	drawPanelCard(auxBounds, battleInsetColor, rl.Fade(battleBorderColor, 0.55))
	medbayText := fmt.Sprintf("Medbay Heal %d/tick", state.Train.MedbayHealingPerTick())
	drawFittedText(medbayText, int32(auxBounds.X+12), int32(auxBounds.Y+10), auxBounds.Width/2-16, 17, 11, rl.Green)
	drawFittedText(lifeSupportText, int32(auxBounds.X+auxBounds.Width/2), int32(auxBounds.Y+10), auxBounds.Width/2-12, 17, 11, lifeSupportColor)

	weaponPanelHeight := float32(46 + len(state.Train.Weapons)*34)
	if weaponPanelHeight < 84 {
		weaponPanelHeight = 84
	}

	combatPanelY := auxBounds.Y + auxBounds.Height + 10
	minCombatPanelHeight := float32(84)
	minGapBetweenPanels := float32(10)
	maxWeaponPanelBottomY := infoPanel.Y + infoPanel.Height - 20
	remainingPanelSpace := maxWeaponPanelBottomY - combatPanelY - minGapBetweenPanels
	if remainingPanelSpace < minCombatPanelHeight+84 {
		remainingPanelSpace = minCombatPanelHeight + 84
	}

	maxPreferredWeaponPanelHeight := remainingPanelSpace - minCombatPanelHeight
	if weaponPanelHeight > maxPreferredWeaponPanelHeight {
		weaponPanelHeight = maxPreferredWeaponPanelHeight
	}
	if weaponPanelHeight < 84 {
		weaponPanelHeight = 84
	}

	weaponPanelY := maxWeaponPanelBottomY - weaponPanelHeight
	combatPanel := rl.Rectangle{
		X:      infoPanel.X + padding,
		Y:      combatPanelY,
		Width:  infoPanel.Width - padding*2,
		Height: weaponPanelY - combatPanelY - minGapBetweenPanels,
	}
	if combatPanel.Height < minCombatPanelHeight {
		combatPanel.Height = minCombatPanelHeight
	}
	b.drawCombatStatusPanel(combatPanel, state)

	weaponPanel := rl.Rectangle{
		X:      infoPanel.X + padding,
		Y:      weaponPanelY,
		Width:  infoPanel.Width - padding*2,
		Height: weaponPanelHeight,
	}
	b.drawWeaponStatusPanel(weaponPanel, state)

	if state.Enemy == nil {
		noEnemyY := int32(combatPanel.Y + combatPanel.Height + 10)
		drawFittedText("No enemy", textX, noEnemyY, infoTextWidth, 24, 14, rl.White)
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
	b.drawEnemySidebar(sidebarBounds, enemySnapshot)
}

func (b BattleWindow) DrawWindowUI(gm *core.GameManager, state *Game) {
}
