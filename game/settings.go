package game

import (
	"fmt"
	"math"

	"github.com/KydZombie/armada/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type SettingsScreen struct {
	previousScreen core.Screen
	selectedRow    int
	draggingSlider int
}

func NewSettingsScreen(currentScreen core.Screen) *SettingsScreen {
	return &SettingsScreen{
		previousScreen: currentScreen,
		selectedRow:    0,
		draggingSlider: -1,
	}
}

func (s *SettingsScreen) ResizeScreen(gm *core.GameManager) {}

func settingsPanelWidth(gm *core.GameManager) float32 {
	width := float32(gm.ScreenWidth) * 0.5
	if width < 360 {
		width = 360
	}
	if width > 480 {
		width = 480
	}
	return width
}

func settingsRowGap(gm *core.GameManager) float32 {
	gap := (float32(gm.ScreenHeight) - 240) / 6
	if gap < 64 {
		gap = 64
	}
	if gap > 88 {
		gap = 88
	}
	return gap
}

func settingsStartY(gm *core.GameManager, gap float32) float32 {
	totalHeight := gap*5 - 110
	centered := float32(gm.ScreenHeight)/2 - totalHeight/2
	minY := float32(180)
	maxY := float32(gm.ScreenHeight) - totalHeight - 72

	if maxY < minY {
		return maxY
	}
	if centered < minY {
		return minY
	}
	if centered > maxY {
		return maxY
	}
	return centered
}

func (s *SettingsScreen) sliderRect(gm *core.GameManager, index int) rl.Rectangle {
	barWidth := settingsPanelWidth(gm)
	barHeight := float32(14)
	rowGap := settingsRowGap(gm)
	startY := settingsStartY(gm, rowGap)
	startX := float32(gm.ScreenWidth)/2 - barWidth/2

	return rl.Rectangle{
		X:      startX,
		Y:      startY + float32(index)*rowGap,
		Width:  barWidth,
		Height: barHeight,
	}
}

func (s *SettingsScreen) toggleRect(gm *core.GameManager, index int) rl.Rectangle {
	barWidth := settingsPanelWidth(gm)
	height := float32(46)
	rowGap := settingsRowGap(gm)
	startY := settingsStartY(gm, rowGap)
	startX := float32(gm.ScreenWidth)/2 - barWidth/2

	return rl.Rectangle{
		X:      startX,
		Y:      startY + float32(3+index)*rowGap,
		Width:  barWidth,
		Height: height,
	}
}

func (s *SettingsScreen) toggleValue(gm *core.GameManager, index int) bool {
	if index == 0 {
		return gm.Fullscreen
	}
	return gm.VSync
}

func (s *SettingsScreen) setToggleValue(gm *core.GameManager, index int, enabled bool) {
	if index == 0 {
		gm.SetFullscreen(enabled)
		return
	}
	gm.SetVSync(enabled)
}

func (s *SettingsScreen) flipToggleValue(gm *core.GameManager, index int) {
	s.setToggleValue(gm, index, !s.toggleValue(gm, index))
}

func clampVolume(v float32) float32 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func quantizeVolume(v float32) float32 {
	return clampVolume(float32(math.Round(float64(v*20)) / 20))
}

func sliderValueFromMouse(rect rl.Rectangle, mouseX float32) float32 {
	return quantizeVolume((mouseX - rect.X) / rect.Width)
}

func (s *SettingsScreen) getVolumeValue(gm *core.GameManager, index int) float32 {
	switch index {
	case 0:
		return gm.MasterVolume
	case 1:
		return gm.MusicVolume
	default:
		return gm.SFXVolume
	}
}

func (s *SettingsScreen) setVolumeValue(gm *core.GameManager, index int, value float32) {
	value = clampVolume(value)
	if index == 0 {
		gm.MasterVolume = value
		return
	}
	if index == 1 {
		gm.MusicVolume = value
		return
	}
	gm.SFXVolume = value
}

func (s *SettingsScreen) UpdateScreen(gm *core.GameManager) {
	volumeLabels := []string{"Master Volume", "Music Volume", "SFX Volume"}
	toggleLabels := []string{"Fullscreen", "VSync"}
	volumeRowCount := len(volumeLabels)
	toggleRowCount := len(toggleLabels)
	backRow := volumeRowCount + toggleRowCount
	maxRow := backRow
	const keyboardVolumeStep = 0.05

	if rl.IsKeyPressed(rl.KeyEscape) {
		gm.SetScreen(s.previousScreen)
		return
	}

	if rl.IsKeyPressed(rl.KeyDown) {
		s.selectedRow = (s.selectedRow + 1) % (maxRow + 1)
	}
	if rl.IsKeyPressed(rl.KeyUp) {
		s.selectedRow = (s.selectedRow - 1 + (maxRow + 1)) % (maxRow + 1)
	}

	if s.selectedRow < volumeRowCount {
		if rl.IsKeyPressed(rl.KeyRight) {
			s.setVolumeValue(gm, s.selectedRow, quantizeVolume(s.getVolumeValue(gm, s.selectedRow)+keyboardVolumeStep))
		}
		if rl.IsKeyPressed(rl.KeyLeft) {
			s.setVolumeValue(gm, s.selectedRow, quantizeVolume(s.getVolumeValue(gm, s.selectedRow)-keyboardVolumeStep))
		}
	}

	if s.selectedRow >= volumeRowCount && s.selectedRow < backRow {
		toggleIndex := s.selectedRow - volumeRowCount
		if rl.IsKeyPressed(rl.KeyLeft) {
			s.setToggleValue(gm, toggleIndex, false)
		}
		if rl.IsKeyPressed(rl.KeyRight) {
			s.setToggleValue(gm, toggleIndex, true)
		}
		if rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeyKpEnter) {
			s.flipToggleValue(gm, toggleIndex)
		}
	}

	if s.selectedRow == backRow && (rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeyKpEnter)) {
		gm.SetScreen(s.previousScreen)
		return
	}

	mousePos := rl.GetMousePosition()
	for i := range volumeLabels {
		sliderRect := s.sliderRect(gm, i)
		hitRect := rl.Rectangle{
			X:      sliderRect.X,
			Y:      sliderRect.Y - 28,
			Width:  sliderRect.Width,
			Height: sliderRect.Height + 56,
		}

		if rl.CheckCollisionPointRec(mousePos, hitRect) {
			s.selectedRow = i
		}

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) && rl.CheckCollisionPointRec(mousePos, hitRect) {
			s.draggingSlider = i
			s.setVolumeValue(gm, i, sliderValueFromMouse(sliderRect, mousePos.X))
		}
	}

	for i := range toggleLabels {
		rect := s.toggleRect(gm, i)
		if rl.CheckCollisionPointRec(mousePos, rect) {
			s.selectedRow = volumeRowCount + i
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				s.flipToggleValue(gm, i)
			}
		}
	}

	if rl.IsMouseButtonDown(rl.MouseLeftButton) && s.draggingSlider >= 0 && s.draggingSlider < volumeRowCount {
		rect := s.sliderRect(gm, s.draggingSlider)
		s.setVolumeValue(gm, s.draggingSlider, sliderValueFromMouse(rect, mousePos.X))
	}
	if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		s.draggingSlider = -1
	}

	rect := rl.Rectangle{
		X:      float32(gm.ScreenWidth)/2 - 170,
		Y:      float32(gm.ScreenHeight) - 110,
		Width:  340,
		Height: 54,
	}
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), rect) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		gm.SetScreen(s.previousScreen)
	}
}

func (s *SettingsScreen) DrawScreen(gm *core.GameManager) {
	rl.DrawTexture(gm.Textures["background"], 0, 0, rl.Color{R: 120, G: 120, B: 120, A: 255})

	rl.DrawTexture(gm.Textures["back"], 0, 0, rl.White)

	rl.DrawCircleGradient(
		0, 0,
		500,
		rl.Color{R: 0, G: 0, B: 0, A: 200},
		rl.Color{R: 0, G: 0, B: 0, A: 0},
	)
	rl.DrawCircleGradient(
		gm.ScreenWidth, 0,
		500,
		rl.Color{R: 0, G: 0, B: 0, A: 200},
		rl.Color{R: 0, G: 0, B: 0, A: 0},
	)
	rl.DrawCircleGradient(
		0, gm.ScreenHeight,
		500,
		rl.Color{R: 0, G: 0, B: 0, A: 200},
		rl.Color{R: 0, G: 0, B: 0, A: 0},
	)
	rl.DrawCircleGradient(
		gm.ScreenWidth, gm.ScreenHeight,
		500,
		rl.Color{R: 0, G: 0, B: 0, A: 200},
		rl.Color{R: 0, G: 0, B: 0, A: 0},
	)

	title := "Settings"

	sizeTitle := rl.MeasureTextEx(gm.Fonts["dh"], title, 100, 2)

	posT := rl.Vector2{
		X: float32(gm.ScreenWidth)/2 - sizeTitle.X/2,
		Y: 70,
	}

	DrawEngravedText(
		gm.Fonts["dh"],
		title,
		posT,
		100,
		2,
		rl.White,
	)

	volumeLabels := []string{"Master Volume", "Music Volume", "SFX Volume"}
	for i, label := range volumeLabels {
		rect := s.sliderRect(gm, i)
		value := s.getVolumeValue(gm, i)
		isSelected := s.selectedRow == i

		labelColor := rl.Black
		if isSelected {
			labelColor = rl.White
		}

		posL := rl.Vector2{
			X: float32(rect.X),
			Y: float32(rect.Y) - 40,
		}

		DrawEngravedText(
			gm.Fonts["dh"],
			label,
			posL,
			30,
			2,
			labelColor,
		)

		rl.DrawRectangleRounded(rect, 0.5, 8, rl.Black)
		fillRect := rl.Rectangle{X: rect.X, Y: rect.Y, Width: rect.Width * value, Height: rect.Height}
		rl.DrawRectangleRounded(fillRect, 0.5, 8, rl.Gray)

		knobX := rect.X + rect.Width*value
		knobRect := rl.Rectangle{X: knobX - 8, Y: rect.Y - 9, Width: 16, Height: rect.Height + 18}
		knobColor := rl.White
		if isSelected {
			knobColor = rl.White
		}
		rl.DrawRectangleRounded(knobRect, 0.4, 6, knobColor)

		percentText := fmt.Sprintf("%d%%", int32(math.Round(float64(value*100))))
		rl.DrawText(percentText, int32(rect.X+rect.Width+20), int32(rect.Y)-2, 24, rl.Color{R: 202, G: 214, B: 233, A: 255})
	}

	toggleLabels := []string{"Fullscreen", "VSync"}
	for i, label := range toggleLabels {
		rect := s.toggleRect(gm, i)
		rowIndex := len(volumeLabels) + i
		isSelected := s.selectedRow == rowIndex
		enabled := s.toggleValue(gm, i)

		rowColor := rl.Black
		textColor := rl.Black
		if isSelected {
			rowColor = rl.Gray
			textColor = rl.White
		}

		rl.DrawRectangleRounded(rect, 0.25, 8, rowColor)

		posL := rl.Vector2{
			X: float32(rect.X + 16),
			Y: float32(rect.Y + 10),
		}

		DrawEngravedText(
			gm.Fonts["dh"],
			label,
			posL,
			30,
			2,
			textColor,
		)

		stateText := "OFF"
		stateColor := rl.Color{R: 194, G: 117, B: 117, A: 255}
		if enabled {
			stateText = "ON"
			stateColor = rl.Color{R: 139, G: 214, B: 159, A: 255}
		}

		stateWidth := rl.MeasureText(stateText, 28)
		rl.DrawText(stateText, int32(rect.X+rect.Width-float32(stateWidth)-18), int32(rect.Y+10), 28, stateColor)
	}

	rect := rl.Rectangle{
		X:      float32(gm.ScreenWidth)/2 - 175,
		Y:      float32(gm.ScreenHeight) - 110,
		Width:  350,
		Height: 90,
	}

	hovered := rl.CheckCollisionPointRec(rl.GetMousePosition(), rect)
	textColor := rl.Black

	if hovered {
		textColor = rl.White
	}

	sizeBack := rl.MeasureTextEx(gm.Fonts["dh"], "Back", 50, 2)

	posB := rl.Vector2{
		X: rect.X + rect.Width/2 - sizeBack.X/2,
		Y: rect.Y + rect.Height/2 - sizeBack.Y/2,
	}

	DrawEngravedText(
		gm.Fonts["dh"],
		"Back",
		posB,
		50,
		2,
		textColor,
	)
}
