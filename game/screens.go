package game

import (
	"fmt"
	"math"

	"github.com/KydZombie/armada/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type MainMenuScreen struct {
	options  []string
	selected int
	hovered  int
}

func NewMainMenuScreen() *MainMenuScreen {
	return &MainMenuScreen{
		options:  []string{"Start Game", "Settings", "Tutorial", "Quit"},
		selected: 0,
		hovered:  -1,
	}
}

func (s *MainMenuScreen) ensureInitialized() {
	if len(s.options) == 0 {
		s.options = []string{"Start Game", "Settings", "Tutorial", "Quit"}
	}
	if s.selected < 0 || s.selected >= len(s.options) {
		s.selected = 0
	}
}

func (s *MainMenuScreen) menuOptionRect(gm *core.GameManager, index int) rl.Rectangle {
	buttonWidth := float32(360)
	buttonHeight := float32(56)
	buttonSpacing := float32(18)
	totalHeight := float32(len(s.options))*buttonHeight + float32(len(s.options)-1)*buttonSpacing
	startY := float32(gm.ScreenHeight)/2 - totalHeight/2 + 40
	x := float32(gm.ScreenWidth)/2 - buttonWidth/2
	y := startY + float32(index)*(buttonHeight+buttonSpacing)

	return rl.Rectangle{
		X:      x,
		Y:      y,
		Width:  buttonWidth,
		Height: buttonHeight,
	}
}

func (s *MainMenuScreen) ResizeScreen(gm *core.GameManager) {}

func (s *MainMenuScreen) UpdateScreen(gm *core.GameManager) {
	s.ensureInitialized()
	optionCount := len(s.options)

	if rl.IsKeyPressed(rl.KeyDown) {
		s.selected = (s.selected + 1) % optionCount
		s.hovered = -1
	}
	if rl.IsKeyPressed(rl.KeyUp) {
		s.selected = (s.selected - 1 + optionCount) % optionCount
		s.hovered = -1
	}

	mousePos := rl.GetMousePosition()
	s.hovered = -1
	for i := range s.options {
		if rl.CheckCollisionPointRec(mousePos, s.menuOptionRect(gm, i)) {
			s.hovered = i
			s.selected = i
			break
		}
	}

	activateSelected := rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeyKpEnter)
	if s.hovered >= 0 && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		activateSelected = true
	}

	if activateSelected {
		s.activateSelected(gm)
	}

	if rl.IsKeyPressed(rl.KeyEscape) {
		gm.Quit()
	}
}

func (s *MainMenuScreen) activateSelected(gm *core.GameManager) {
	switch s.options[s.selected] {
	case "Start Game":
		gm.SetScreen(NewGameScreen(gm))
	case "Settings":
		gm.SetScreen(NewSettingsScreen(s))
	case "Tutorial":
		gm.SetScreen(NewTutorialScreen(s))
	case "Quit":
		gm.Quit()
	}
}

func (s *MainMenuScreen) DrawScreen(gm *core.GameManager) {
	rl.ClearBackground(rl.Color{R: 10, G: 18, B: 34, A: 255})

	centerX := gm.ScreenWidth / 2
	title := "ARMADA"
	titleWidth := rl.MeasureText(title, 92)
	rl.DrawText(title, centerX-titleWidth/2, 88, 92, rl.Color{R: 229, G: 237, B: 255, A: 255})

	subtitle := "Select an option"
	subtitleWidth := rl.MeasureText(subtitle, 30)
	rl.DrawText(subtitle, centerX-subtitleWidth/2, 190, 30, rl.Color{R: 150, G: 177, B: 210, A: 255})

	for i, option := range s.options {
		rect := s.menuOptionRect(gm, i)
		isSelected := s.selected == i
		isHovered := s.hovered == i

		fillColor := rl.Color{R: 28, G: 48, B: 77, A: 255}
		textColor := rl.Color{R: 204, G: 217, B: 236, A: 255}

		if isHovered {
			fillColor = rl.Color{R: 42, G: 77, B: 119, A: 255}
		}
		if isSelected {
			fillColor = rl.Color{R: 68, G: 118, B: 171, A: 255}
			textColor = rl.White
		}

		rl.DrawRectangleRec(rect, fillColor)
		rl.DrawRectangleLinesEx(rect, 2, rl.Color{R: 102, G: 137, B: 179, A: 255})

		textSize := int32(30)
		textWidth := rl.MeasureText(option, textSize)
		textX := int32(rect.X + rect.Width/2 - float32(textWidth)/2)
		textY := int32(rect.Y + rect.Height/2 - float32(textSize)/2)
		rl.DrawText(option, textX, textY, textSize, textColor)
	}
}

func (s *MainMenuScreen) DrawScreenUI(gm *core.GameManager) {
	controls := "Keyboard: Up/Down + Enter | Mouse: Hover + Click | Esc: Quit"
	controlsWidth := rl.MeasureText(controls, 20)
	rl.DrawText(
		controls,
		gm.ScreenWidth/2-controlsWidth/2,
		gm.ScreenHeight-48,
		20,
		rl.Color{R: 140, G: 160, B: 188, A: 255},
	)
}

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
	totalHeight := gap*5 + 54
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

func (s *SettingsScreen) backButtonRect(gm *core.GameManager) rl.Rectangle {
	barWidth := settingsPanelWidth(gm)
	buttonWidth := barWidth * 0.8
	if buttonWidth < 280 {
		buttonWidth = 280
	}
	if buttonWidth > 340 {
		buttonWidth = 340
	}
	gap := settingsRowGap(gm)
	startY := settingsStartY(gm, gap)
	return rl.Rectangle{
		X:      float32(gm.ScreenWidth)/2 - buttonWidth/2,
		Y:      startY + gap*5,
		Width:  buttonWidth,
		Height: 54,
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

	backRect := s.backButtonRect(gm)
	if rl.CheckCollisionPointRec(mousePos, backRect) {
		s.selectedRow = backRow
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			gm.SetScreen(s.previousScreen)
			return
		}
	}
}

func (s *SettingsScreen) DrawScreen(gm *core.GameManager) {
	rl.ClearBackground(rl.Color{R: 14, G: 21, B: 37, A: 255})

	title := "Settings"
	titleWidth := rl.MeasureText(title, 64)
	rl.DrawText(title, gm.ScreenWidth/2-titleWidth/2, 72, 64, rl.Color{R: 230, G: 236, B: 247, A: 255})

	volumeLabels := []string{"Master Volume", "Music Volume", "SFX Volume"}
	for i, label := range volumeLabels {
		rect := s.sliderRect(gm, i)
		value := s.getVolumeValue(gm, i)
		isSelected := s.selectedRow == i

		labelColor := rl.Color{R: 186, G: 199, B: 220, A: 255}
		if isSelected {
			labelColor = rl.White
		}
		rl.DrawText(label, int32(rect.X), int32(rect.Y)-40, 30, labelColor)

		rl.DrawRectangleRounded(rect, 0.5, 8, rl.Color{R: 37, G: 58, B: 87, A: 255})
		fillRect := rl.Rectangle{X: rect.X, Y: rect.Y, Width: rect.Width * value, Height: rect.Height}
		rl.DrawRectangleRounded(fillRect, 0.5, 8, rl.Color{R: 87, G: 150, B: 221, A: 255})

		knobX := rect.X + rect.Width*value
		knobRect := rl.Rectangle{X: knobX - 8, Y: rect.Y - 9, Width: 16, Height: rect.Height + 18}
		knobColor := rl.Color{R: 220, G: 230, B: 244, A: 255}
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

		rowColor := rl.Color{R: 31, G: 49, B: 73, A: 255}
		textColor := rl.Color{R: 197, G: 209, B: 227, A: 255}
		if isSelected {
			rowColor = rl.Color{R: 60, G: 94, B: 135, A: 255}
			textColor = rl.White
		}

		rl.DrawRectangleRounded(rect, 0.25, 8, rowColor)
		rl.DrawRectangleLinesEx(rect, 2, rl.Color{R: 102, G: 137, B: 179, A: 255})
		rl.DrawText(label, int32(rect.X+16), int32(rect.Y+10), 28, textColor)

		stateText := "OFF"
		stateColor := rl.Color{R: 194, G: 117, B: 117, A: 255}
		if enabled {
			stateText = "ON"
			stateColor = rl.Color{R: 139, G: 214, B: 159, A: 255}
		}

		stateWidth := rl.MeasureText(stateText, 28)
		rl.DrawText(stateText, int32(rect.X+rect.Width-float32(stateWidth)-18), int32(rect.Y+10), 28, stateColor)
	}

	backRect := s.backButtonRect(gm)
	backColor := rl.Color{R: 34, G: 53, B: 79, A: 255}
	backTextColor := rl.Color{R: 197, G: 209, B: 227, A: 255}
	if s.selectedRow == len(volumeLabels)+len(toggleLabels) {
		backColor = rl.Color{R: 65, G: 100, B: 143, A: 255}
		backTextColor = rl.White
	}

	rl.DrawRectangleRec(backRect, backColor)
	rl.DrawRectangleLinesEx(backRect, 2, rl.Color{R: 102, G: 137, B: 179, A: 255})
	backText := "Back"
	backWidth := rl.MeasureText(backText, 30)
	rl.DrawText(backText, int32(backRect.X+backRect.Width/2-float32(backWidth)/2), int32(backRect.Y+backRect.Height/2-15), 30, backTextColor)
}

func (s *SettingsScreen) DrawScreenUI(gm *core.GameManager) {
	controls := "Keyboard: Up/Down + Left/Right + Enter | Mouse: Drag sliders + Click toggles/back | Esc: Back"
	controlsWidth := rl.MeasureText(controls, 20)
	rl.DrawText(
		controls,
		gm.ScreenWidth/2-controlsWidth/2,
		gm.ScreenHeight-44,
		20,
		rl.Color{R: 140, G: 160, B: 188, A: 255},
	)
}

type TutorialScreen struct {
	previousScreen core.Screen
}

func NewTutorialScreen(currentScreen core.Screen) *TutorialScreen {
	return &TutorialScreen{previousScreen: currentScreen}
}

func (s *TutorialScreen) ResizeScreen(gm *core.GameManager) {}

func (s *TutorialScreen) UpdateScreen(gm *core.GameManager) {
	if rl.IsKeyPressed(rl.KeyEscape) || rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeyKpEnter) {
		gm.SetScreen(s.previousScreen)
		return
	}

	backRect := rl.Rectangle{
		X:      float32(gm.ScreenWidth)/2 - 170,
		Y:      float32(gm.ScreenHeight) - 110,
		Width:  340,
		Height: 54,
	}
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), backRect) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		gm.SetScreen(s.previousScreen)
	}
}

func (s *TutorialScreen) DrawScreen(gm *core.GameManager) {
	rl.ClearBackground(rl.Color{R: 11, G: 18, B: 31, A: 255})

	title := "Tutorial"
	titleWidth := rl.MeasureText(title, 64)
	rl.DrawText(title, gm.ScreenWidth/2-titleWidth/2, 70, 64, rl.Color{R: 230, G: 236, B: 247, A: 255})

	lines := []string{
		"1. Keep your ship running while under pressure.",
		"2. Enter commands in the terminal window during gameplay.",
		"3. Watch enemy and train windows to react quickly.",
		"4. Use command 'help' in-game to learn available actions.",
	}

	startY := int32(190)
	for i, line := range lines {
		rl.DrawText(line, 120, startY+int32(i*48), 34, rl.Color{R: 191, G: 205, B: 227, A: 255})
	}

	backRect := rl.Rectangle{
		X:      float32(gm.ScreenWidth)/2 - 170,
		Y:      float32(gm.ScreenHeight) - 110,
		Width:  340,
		Height: 54,
	}

	hovered := rl.CheckCollisionPointRec(rl.GetMousePosition(), backRect)
	backColor := rl.Color{R: 35, G: 54, B: 82, A: 255}
	if hovered {
		backColor = rl.Color{R: 62, G: 94, B: 137, A: 255}
	}
	textColor := rl.Color{R: 204, G: 217, B: 236, A: 255}
	if hovered {
		textColor = rl.White
	}

	rl.DrawRectangleRec(backRect, backColor)
	rl.DrawRectangleLinesEx(backRect, 2, rl.Color{R: 102, G: 137, B: 179, A: 255})
	backText := "Back"
	backWidth := rl.MeasureText(backText, 30)
	rl.DrawText(backText, int32(backRect.X+backRect.Width/2-float32(backWidth)/2), int32(backRect.Y+backRect.Height/2-15), 30, textColor)
}

func (s *TutorialScreen) DrawScreenUI(gm *core.GameManager) {
	controls := "Press Enter/Esc or click Back to return"
	controlsWidth := rl.MeasureText(controls, 22)
	rl.DrawText(controls, gm.ScreenWidth/2-controlsWidth/2, gm.ScreenHeight-42, 22, rl.Color{R: 140, G: 160, B: 188, A: 255})
}
