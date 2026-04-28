package game

import (
	"fmt"
	"strings"

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
	return false
}

func (m MissionWindow) UpdateWindow(gm *core.GameManager, state *Game) {
}

func (m MissionWindow) DrawWindow(gm *core.GameManager, state *Game) {
	bounds := m.GetBounds()
	rl.DrawRectangleRec(bounds, battleBgColor)
	rl.DrawRectangleLinesEx(bounds, 2, rl.Fade(battleBorderColor, 0.6))

	const missionFontSize int32 = 17
	const missionFontMinSize int32 = 12
	const missionTitleFontSize int32 = 26
	const missionTitleMinSize int32 = 18
	const missionLeadFontSize int32 = 22
	const missionLeadMinSize int32 = 15

	padding := float32(12)
	content := rl.Rectangle{
		X:      bounds.X + padding,
		Y:      bounds.Y + padding,
		Width:  bounds.Width - padding*2,
		Height: bounds.Height - padding*2,
	}
	drawPanelCard(content, battlePanelColor, rl.Fade(battleBorderColor, 0.85))

	header := rl.Rectangle{X: content.X + 10, Y: content.Y + 10, Width: content.Width - 20, Height: 34}
	drawPanelCard(header, battleInsetColor, rl.Fade(battleAccentColor, 0.8))
	rl.DrawRectangleRec(rl.Rectangle{X: header.X, Y: header.Y, Width: 4, Height: header.Height}, battleAccentColor)
	drawFittedText("Mission", int32(header.X+12), int32(header.Y+6), header.Width-24, missionTitleFontSize, missionTitleMinSize, rl.White)

	missionBounds := rl.Rectangle{X: content.X + 10, Y: header.Y + header.Height + 8, Width: content.Width - 20, Height: 136}
	drawPanelCard(missionBounds, battleInsetColor, rl.Fade(battleBorderColor, 0.55))

	progressLine := fmt.Sprintf("Kill hostiles: %d/%d", state.Wave.KillsDone, state.Wave.KillsRequired)
	progressColor := rl.SkyBlue
	if state.Wave.Success {
		progressColor = rl.Green
	} else if state.Wave.Failed {
		progressColor = rl.Red
	}
	drawFittedText(progressLine, int32(missionBounds.X+12), int32(missionBounds.Y+10), missionBounds.Width-24, missionLeadFontSize, missionLeadMinSize, progressColor)

	timeLine := fmt.Sprintf("Time left: %.0fs", state.Wave.TimeRemaining)
	timeColor := rl.SkyBlue
	if state.Wave.TimeRemaining <= 15 {
		timeColor = rl.Yellow
	}
	if state.Wave.TimeRemaining <= 8 {
		timeColor = rl.Orange
	}
	if state.Wave.TimeRemaining <= 4 {
		timeColor = rl.Red
	}
	drawFittedText(timeLine, int32(missionBounds.X+12), int32(missionBounds.Y+36), missionBounds.Width-24, missionLeadFontSize, missionLeadMinSize, timeColor)

	targetLine := fmt.Sprintf("Target zone: %s", strings.ToUpper(string(state.SelectedTargetZone)))
	drawFittedText(targetLine, int32(missionBounds.X+12), int32(missionBounds.Y+58), missionBounds.Width-24, missionFontSize, missionFontMinSize, rl.Gold)

	rulesLine := "Win: kills target reached. Lose: timer 0 or hull 0."
	drawFittedText(rulesLine, int32(missionBounds.X+12), int32(missionBounds.Y+80), missionBounds.Width-24, missionFontSize, missionFontMinSize, rl.LightGray)

	drawFittedText("Crew positions", int32(missionBounds.X+12), int32(missionBounds.Y+100), missionBounds.Width-24, missionFontSize, missionFontMinSize, rl.White)
	crewLine := state.CrewSupportSummaryText()
	crewColor := rl.Green
	if crewLine == "none" {
		crewColor = rl.LightGray
	}
	drawFittedText(crewLine, int32(missionBounds.X+12), int32(missionBounds.Y+118), missionBounds.Width-24, missionFontSize, missionFontMinSize, crewColor)

	drawFittedText("Cart damage", int32(missionBounds.X+12), int32(missionBounds.Y+138), missionBounds.Width-24, missionFontSize, missionFontMinSize, rl.White)
	threatLine := state.ThreatSummaryText()
	threatColor := rl.Orange
	if threatLine == "none" {
		threatColor = rl.LightGray
	}
	drawFittedText(threatLine, int32(missionBounds.X+12), int32(missionBounds.Y+156), missionBounds.Width-24, missionFontSize, missionFontMinSize, threatColor)

	if state.Wave.Success || state.Wave.Failed {
		bannerText := "LEVEL FAILED"
		bannerColor := rl.Red
		bannerFill := rl.NewColor(90, 20, 22, 245)
		if state.Wave.Success {
			bannerText = "LEVEL CLEARED"
			bannerColor = rl.Green
			bannerFill = rl.NewColor(18, 72, 36, 245)
		}

		bannerBounds := rl.Rectangle{X: content.X + 10, Y: missionBounds.Y + missionBounds.Height + 8, Width: content.Width - 20, Height: 48}
		drawPanelCard(bannerBounds, bannerFill, rl.Fade(bannerColor, 0.9))
		drawFittedText(bannerText, int32(bannerBounds.X+12), int32(bannerBounds.Y+11), bannerBounds.Width-24, missionLeadFontSize, missionLeadMinSize, bannerColor)
	}

	logBounds := rl.Rectangle{X: content.X + 10, Y: missionBounds.Y + missionBounds.Height + 60, Width: content.Width - 20, Height: content.Y + content.Height - (missionBounds.Y + missionBounds.Height + 70)}
	if logBounds.Height >= 70 {
		drawPanelCard(logBounds, rl.Fade(battleInsetColor, 0.92), rl.Fade(battleBorderColor, 0.45))
		drawFittedText("Action Log", int32(logBounds.X+12), int32(logBounds.Y+8), logBounds.Width-24, missionFontSize, missionFontMinSize, rl.White)
		lineY := int32(logBounds.Y + 34)
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
}
