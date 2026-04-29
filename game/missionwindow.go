package game

import (
	"fmt"

	"github.com/KydZombie/armada/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type MissionWindow struct {
	core.BaseWindow[Game]
}

func NewMissionWindow(sizeFunc func(gm *core.GameManager) rl.Rectangle, gm *core.GameManager) *MissionWindow {
	return &MissionWindow{
		BaseWindow: core.NewBaseWindow[Game](sizeFunc, gm, true),
	}
}

func (m MissionWindow) HandleInput(gm *core.GameManager, state *Game) bool {
	if state.isMissionBriefingActive() {
		mousePos := rl.GetMousePosition()
		if !rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			return false
		}

		if rl.CheckCollisionPointRec(mousePos, missionBriefingStartButtonRect(gm)) {
			state.MissionBriefingActive = false
			return true
		}

		return true
	}

	if state.isGameOverModalActive() {
		mousePos := rl.GetMousePosition()
		if !rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			return false
		}

		if rl.CheckCollisionPointRec(mousePos, gameOverRetryButtonRect(gm)) {
			gm.SetScreen(NewGameScreen(gm))
			return true
		}
		if rl.CheckCollisionPointRec(mousePos, gameOverMenuButtonRect(gm)) {
			gm.SetScreen(NewMainMenuScreen())
			return true
		}

		return true
	}

	return false
}

func (m MissionWindow) UpdateWindow(gm *core.GameManager, state *Game) {
}

func (m MissionWindow) DrawWindow(gm *core.GameManager, state *Game) {
	bounds := m.GetBounds()

	const missionFontSize int32 = 21
	const missionFontMinSize int32 = 15
	const missionTitleFontSize int32 = 30
	const missionTitleMinSize int32 = 22
	const missionLeadFontSize int32 = 26
	const missionLeadMinSize int32 = 18

	padding := float32(12)
	content := rl.Rectangle{
		X:      bounds.X + padding,
		Y:      bounds.Y + padding,
		Width:  bounds.Width - padding*2,
		Height: bounds.Height - padding*2,
	}
	drawPanelCard(content, battlePanelColor, rl.Fade(battleBorderColor, 0.85))

	header := rl.Rectangle{X: content.X + 10, Y: content.Y + 10, Width: content.Width - 20, Height: 34}
	header.Height = 38
	drawPanelCard(header, battleInsetColor, rl.Fade(battleAccentColor, 0.8))
	rl.DrawRectangleRec(rl.Rectangle{X: header.X, Y: header.Y, Width: 4, Height: header.Height}, battleAccentColor)
	drawFittedText(gm, "Mission", int32(header.X+12), int32(header.Y+6), header.Width-24, missionTitleFontSize, missionTitleMinSize, rl.White)

	missionBounds := rl.Rectangle{X: content.X + 10, Y: header.Y + header.Height + 10, Width: content.Width - 20, Height: 140}
	drawPanelCard(missionBounds, battleInsetColor, rl.Fade(battleBorderColor, 0.55))

	progressLine := fmt.Sprintf("Kill hostiles: %d/%d", state.Wave.KillsDone, state.Wave.KillsRequired)
	progressColor := rl.LightGray
	if state.Wave.Success {
		progressColor = rl.Green
	} else if state.Wave.Failed {
		progressColor = rl.Red
	}
	drawFittedText(gm, progressLine, int32(missionBounds.X+12), int32(missionBounds.Y+10), missionBounds.Width-24, missionLeadFontSize, missionLeadMinSize, progressColor)

	timeLine := fmt.Sprintf("Time left: %.0fs", state.Wave.TimeRemaining)
	timeColor := rl.Gray
	if state.Wave.TimeRemaining <= 15 {
		timeColor = rl.Yellow
	}
	if state.Wave.TimeRemaining <= 8 {
		timeColor = rl.Orange
	}
	if state.Wave.TimeRemaining <= 4 {
		timeColor = rl.Red
	}
	drawFittedText(gm, timeLine, int32(missionBounds.X+12), int32(missionBounds.Y+40), missionBounds.Width-24, missionLeadFontSize, missionLeadMinSize, timeColor)

	if state.Wave.Success {
		bannerText := "LEVEL FAILED"
		bannerColor := rl.Green
		bannerFill := rl.NewColor(18, 72, 36, 245)
		bannerText = "LEVEL CLEARED"

		bannerBounds := rl.Rectangle{X: content.X + 10, Y: missionBounds.Y + missionBounds.Height + 8, Width: content.Width - 20, Height: 48}
		drawPanelCard(bannerBounds, bannerFill, rl.Fade(bannerColor, 0.9))
		drawFittedText(gm, bannerText, int32(bannerBounds.X+12), int32(bannerBounds.Y+11), bannerBounds.Width-24, missionLeadFontSize, missionLeadMinSize, bannerColor)
	}

	logBounds := rl.Rectangle{X: content.X + 10, Y: missionBounds.Y + missionBounds.Height + 60, Width: content.Width - 20, Height: content.Y + content.Height - (missionBounds.Y + missionBounds.Height + 70)}
	if logBounds.Height >= 70 {
		drawPanelCard(logBounds, rl.Fade(battleInsetColor, 0.92), rl.Fade(battleBorderColor, 0.45))
		drawFittedText(gm, "Action Log", int32(logBounds.X+12), int32(logBounds.Y+8), logBounds.Width-24, missionFontSize, missionFontMinSize, rl.White)
		lineY := int32(logBounds.Y + 36)
		maxY := int32(logBounds.Y + logBounds.Height - 10)
		for i, line := range state.combatStatusLines {
			if i >= 3 {
				break
			}
			size := fitTextSize(line, logBounds.Width-34, missionFontSize, missionFontMinSize)
			wrapped := wrapTerminalLine(line, logBounds.Width-34, size)
			for j, wrappedLine := range wrapped {
				if lineY+size > maxY {
					break
				}
				if j == 0 {
					rl.DrawCircle(int32(logBounds.X+18), lineY+size/2, 3, rl.LightGray)
				}
				rl.DrawText(wrappedLine, int32(logBounds.X+28), lineY, size, rl.LightGray)
				lineY += size + 2
			}
			lineY += 5
		}
	}
}

func (m MissionWindow) DrawWindowUI(gm *core.GameManager, state *Game) {
	if state.isMissionBriefingActive() {
		drawMissionBriefingModal(gm)
		return
	}

	if !state.isGameOverModalActive() {
		return
	}

	drawGameOverModal(gm)
}

func drawMissionListSection(gm *core.GameManager, bounds rl.Rectangle, title string, items []string, accent rl.Color, titleFontSize int32, itemFontSize int32) {
	drawPanelCard(bounds, battleInsetColor, rl.Fade(accent, 0.65))
	rl.DrawRectangleRec(rl.Rectangle{X: bounds.X, Y: bounds.Y, Width: 4, Height: bounds.Height}, accent)
	drawFittedText(gm, title, int32(bounds.X+10), int32(bounds.Y+6), bounds.Width-20, titleFontSize, titleFontSize-4, rl.White)

	lineY := int32(bounds.Y + 34)
	maxY := int32(bounds.Y + bounds.Height - 8)
	if len(items) == 0 {
		items = []string{"none"}
	}
	for _, item := range items {
		size := fitTextSize(item, bounds.Width-20, itemFontSize, itemFontSize-4)
		wrapped := wrapTerminalLine(item, bounds.Width-20, size)
		for _, line := range wrapped {
			if lineY+size > maxY {
				return
			}
			rl.DrawText(line, int32(bounds.X+10), lineY, size, accent)
			lineY += size + 2
		}
		lineY += 3
	}
}

func drawMissionBriefingModal(gm *core.GameManager) {
	overlay := rl.Rectangle{X: 0, Y: 0, Width: float32(gm.ScreenWidth), Height: float32(gm.ScreenHeight)}
	rl.DrawRectangleRec(overlay, rl.NewColor(4, 8, 16, 200))

	panelWidth := float32(gm.ScreenWidth) * 0.52
	if panelWidth < 480 {
		panelWidth = 480
	}
	if panelWidth > 760 {
		panelWidth = 760
	}
	panelHeight := float32(250)
	panelBounds := rl.Rectangle{
		X:      float32(gm.ScreenWidth)/2 - panelWidth/2,
		Y:      float32(gm.ScreenHeight)/2 - panelHeight/2,
		Width:  panelWidth,
		Height: panelHeight,
	}
	drawPanelCard(panelBounds, rl.Black, rl.Gray)

	title := "HOW TO WIN"
	titleSize := fitTextSize(title, panelBounds.Width-48, 44, 28)
	titleWidth := int32(rl.MeasureTextEx(gm.Fonts["ec-i"], title, float32(titleSize), 2).X)
	rl.DrawTextEx(
		gm.Fonts["ec-i"],
		title,
		rl.NewVector2(
			float32(panelBounds.X+panelBounds.Width/2-float32(titleWidth)/2),
			float32(panelBounds.Y+18),
		),
		float32(titleSize),
		2,
		rl.White,
	)

	message := "Kill enough hostiles before the timer runs out. Keep the train hull alive."
	messageSize := fitTextSize(message, panelBounds.Width-48, 30, 20)
	messageLines := wrapTerminalLine(message, panelBounds.Width-48, messageSize)
	lineY := int32(panelBounds.Y + 80)
	for _, line := range messageLines {
		lineWidth := int32(rl.MeasureTextEx(gm.Fonts["ec-i"], line, float32(messageSize), 2).X)
		rl.DrawTextEx(
			gm.Fonts["ec-i"],
			line,
			rl.NewVector2(
				float32(panelBounds.X+panelBounds.Width/2-float32(lineWidth)/2),
				float32(lineY),
			),
			float32(messageSize),
			2,
			rl.LightGray,
		)
		lineY += messageSize + 5
	}

	startRect := missionBriefingStartButtonRect(gm)
	buttonColor := rl.DarkGray
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), startRect) {
		buttonColor = rl.Gray
	}
	drawPanelCard(startRect, buttonColor, rl.DarkGray)
	buttonText := "Start"
	buttonTextSize := fitTextSize(buttonText, startRect.Width-24, 28, 18)
	buttonWidth := int32(rl.MeasureTextEx(gm.Fonts["ec-i"], buttonText, float32(buttonTextSize), 2).X)
	rl.DrawTextEx(
		gm.Fonts["ec-i"],
		buttonText,
		rl.NewVector2(
			float32(startRect.X+startRect.Width/2-float32(buttonWidth)/2),
			float32(startRect.Y+startRect.Height/2-float32(buttonTextSize)/2),
		),
		float32(buttonTextSize),
		2,
		rl.White,
	)
}

func drawGameOverModal(gm *core.GameManager) {
	overlay := rl.Rectangle{X: 0, Y: 0, Width: float32(gm.ScreenWidth), Height: float32(gm.ScreenHeight)}
	rl.DrawRectangleRec(overlay, rl.NewColor(4, 8, 16, 200))

	panelWidth := float32(gm.ScreenWidth) * 0.5
	if panelWidth < 460 {
		panelWidth = 460
	}
	if panelWidth > 760 {
		panelWidth = 760
	}
	panelHeight := float32(200)
	panelBounds := rl.Rectangle{
		X:      float32(gm.ScreenWidth)/2 - panelWidth/2,
		Y:      float32(gm.ScreenHeight)/2 - panelHeight/2,
		Width:  panelWidth,
		Height: panelHeight,
	}
	drawPanelCard(panelBounds, rl.Black, rl.Gray)
	rl.DrawRectangleLinesEx(panelBounds, 4, rl.DarkGray)

	title := "LEVEL FAILED"
	titleSize := fitTextSize(title, panelBounds.Width-48, 36, 24)
	titleWidth := int32(rl.MeasureTextEx(gm.Fonts["ec-i"], title, float32(titleSize), 2).X)
	rl.DrawTextEx(
		gm.Fonts["ec-i"],
		title,
		rl.NewVector2(
			float32(panelBounds.X+panelBounds.Width/2-float32(titleWidth)/2),
			float32(panelBounds.Y+24),
		),
		float32(titleSize),
		2,
		rl.White,
	)

	retryRect := gameOverRetryButtonRect(gm)
	menuRect := gameOverMenuButtonRect(gm)
	mousePos := rl.GetMousePosition()
	retryHover := rl.CheckCollisionPointRec(mousePos, retryRect)
	menuHover := rl.CheckCollisionPointRec(mousePos, menuRect)

	retryFill := rl.DarkGray
	if retryHover {
		retryFill = rl.Gray
	}
	menuFill := rl.DarkGray
	if menuHover {
		menuFill = rl.Gray
	}

	drawPanelCard(retryRect, retryFill, rl.DarkGray)
	drawPanelCard(menuRect, menuFill, rl.DarkGray)

	retryText := "Retry"
	menuText := "Main Menu"
	retryTextSize := fitTextSize(retryText, retryRect.Width-24, 24, 16)
	menuTextSize := fitTextSize(menuText, menuRect.Width-24, 24, 16)
	retryWidth := int32(rl.MeasureTextEx(gm.Fonts["ec-i"], retryText, float32(retryTextSize), 2).X)
	menuWidth := int32(rl.MeasureTextEx(gm.Fonts["ec-i"], menuText, float32(menuTextSize), 2).X)
	rl.DrawTextEx(
		gm.Fonts["ec-i"],
		retryText,
		rl.NewVector2(
			float32(retryRect.X+retryRect.Width/2-float32(retryWidth)/2),
			float32(retryRect.Y+retryRect.Height/2-float32(retryTextSize)/2),
		),
		float32(retryTextSize),
		2,
		rl.White,
	)
	rl.DrawTextEx(
		gm.Fonts["ec-i"],
		menuText,
		rl.NewVector2(
			float32(menuRect.X+menuRect.Width/2-float32(menuWidth)/2),
			float32(menuRect.Y+menuRect.Height/2-float32(menuTextSize)/2),
		),
		float32(menuTextSize),
		2,
		rl.White,
	)
}

func missionBriefingStartButtonRect(gm *core.GameManager) rl.Rectangle {
	buttonWidth := float32(220)
	buttonHeight := float32(58)
	return rl.Rectangle{
		X:      float32(gm.ScreenWidth)/2 - buttonWidth/2,
		Y:      float32(gm.ScreenHeight)/2 + 56,
		Width:  buttonWidth,
		Height: buttonHeight,
	}
}

func gameOverRetryButtonRect(gm *core.GameManager) rl.Rectangle {
	buttonWidth := float32(180)
	buttonHeight := float32(54)
	gap := float32(18)
	totalWidth := buttonWidth*2 + gap
	x := float32(gm.ScreenWidth)/2 - totalWidth/2
	y := float32(gm.ScreenHeight)/2 + 70
	return rl.Rectangle{X: x, Y: y, Width: buttonWidth, Height: buttonHeight}
}

func gameOverMenuButtonRect(gm *core.GameManager) rl.Rectangle {
	buttonWidth := float32(180)
	buttonHeight := float32(54)
	gap := float32(18)
	totalWidth := buttonWidth*2 + gap
	x := float32(gm.ScreenWidth)/2 - totalWidth/2 + buttonWidth + gap
	y := float32(gm.ScreenHeight)/2 + 70
	return rl.Rectangle{X: x, Y: y, Width: buttonWidth, Height: buttonHeight}
}
