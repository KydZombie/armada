package game

import (
	"fmt"

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

func (s *SettingsScreen) sliderRect(gm *core.GameManager, index int) rl.Rectangle {
	barWidth := float32(420)
	barHeight := float32(14)
	startX := float32(gm.ScreenWidth)/2 - barWidth/2
	startY := float32(gm.ScreenHeight)/2 - 110
	rowGap := float32(92)

	return rl.Rectangle{
		X:      startX,
		Y:      startY + float32(index)*rowGap,
		Width:  barWidth,
		Height: barHeight,
	}
}

func (s *SettingsScreen) backButtonRect(gm *core.GameManager) rl.Rectangle {
	return rl.Rectangle{
		X:      float32(gm.ScreenWidth)/2 - 170,
		Y:      float32(gm.ScreenHeight) - 130,
		Width:  340,
		Height: 54,
	}
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

func sliderValueFromMouse(rect rl.Rectangle, mouseX float32) float32 {
	return clampVolume((mouseX - rect.X) / rect.Width)
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
	labels := []string{"Master Volume", "Music Volume", "SFX Volume"}
	maxRow := len(labels)

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

	if s.selectedRow < maxRow {
		if rl.IsKeyDown(rl.KeyRight) {
			s.setVolumeValue(gm, s.selectedRow, s.getVolumeValue(gm, s.selectedRow)+0.01)
		}
		if rl.IsKeyDown(rl.KeyLeft) {
			s.setVolumeValue(gm, s.selectedRow, s.getVolumeValue(gm, s.selectedRow)-0.01)
		}
	}

	if s.selectedRow == maxRow && (rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeyKpEnter)) {
		gm.SetScreen(s.previousScreen)
		return
	}

	mousePos := rl.GetMousePosition()
	for i := range labels {
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

	if rl.IsMouseButtonDown(rl.MouseLeftButton) && s.draggingSlider >= 0 && s.draggingSlider < maxRow {
		rect := s.sliderRect(gm, s.draggingSlider)
		s.setVolumeValue(gm, s.draggingSlider, sliderValueFromMouse(rect, mousePos.X))
	}
	if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		s.draggingSlider = -1
	}

	backRect := s.backButtonRect(gm)
	if rl.CheckCollisionPointRec(mousePos, backRect) {
		s.selectedRow = maxRow
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

	labels := []string{"Master Volume", "Music Volume", "SFX Volume"}
	for i, label := range labels {
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

		percentText := fmt.Sprintf("%d%%", int32(value*100))
		rl.DrawText(percentText, int32(rect.X+rect.Width+20), int32(rect.Y)-2, 24, rl.Color{R: 202, G: 214, B: 233, A: 255})
	}

	backRect := s.backButtonRect(gm)
	backColor := rl.Color{R: 34, G: 53, B: 79, A: 255}
	backTextColor := rl.Color{R: 197, G: 209, B: 227, A: 255}
	if s.selectedRow == len(labels) {
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
	controls := "Keyboard: Up/Down + Left/Right + Enter | Mouse: Drag sliders + Click Back | Esc: Back"
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
