package game

import (
	"github.com/KydZombie/armada/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func DrawEngravedText(font rl.Font, text string, pos rl.Vector2, size float32, spacing float32, color rl.Color) {
	shadowOffset := float32(2)

	rl.DrawTextEx(
		font,
		text,
		rl.Vector2{X: pos.X - shadowOffset, Y: pos.Y - shadowOffset},
		size,
		spacing,
		rl.Color{R: 0, G: 0, B: 0, A: 180},
	)

	rl.DrawTextEx(
		font,
		text,
		rl.Vector2{X: pos.X + shadowOffset, Y: pos.Y + shadowOffset},
		size,
		spacing,
		rl.Color{R: 255, G: 255, B: 255, A: 60},
	)

	rl.DrawTextEx(
		font,
		text,
		pos,
		size,
		spacing,
		color,
	)
}

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
	buttonWidth := float32(380)
	buttonHeight := float32(95)
	buttonSpacing := float32(12.5)
	totalHeight := float32(len(s.options))*buttonHeight + float32(len(s.options)-1)*buttonSpacing
	startY := float32(gm.ScreenHeight)/2 - totalHeight/2 + 137.5
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
	rl.DrawTexture(gm.Textures["background"], 0, 0, rl.Color{R: 120, G: 120, B: 120, A: 255})

	rl.DrawTexture(gm.Textures["options"], 0, 0, rl.White)

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

	rl.DrawTexture(gm.Textures["title"], 0, 0, rl.White)

	for i, option := range s.options {
		rect := s.menuOptionRect(gm, i)
		isSelected := s.selected == i
		isHovered := s.hovered == i

		fillColor := rl.Color{R: 0, G: 0, B: 0, A: 0}
		textColor := rl.Color{R: 0, G: 0, B: 0, A: 255}

		if isHovered {
			fillColor = rl.Color{R: 0, G: 0, B: 0, A: 0}
		}
		if isSelected {
			fillColor = rl.Color{R: 0, G: 0, B: 0, A: 0}
			textColor = rl.White
		}

		rl.DrawRectangleRec(rect, fillColor)

		sizeOption := rl.MeasureTextEx(gm.Fonts["dh"], option, 50, 2)

		pos := rl.Vector2{
			X: rect.X + rect.Width/2 - sizeOption.X/2,
			Y: rect.Y + rect.Height/2 - sizeOption.Y/2,
		}

		DrawEngravedText(
			gm.Fonts["dh"],
			option,
			pos,
			50,
			2,
			textColor,
		)
	}
}

func (s *MainMenuScreen) DrawScreenUI(gm *core.GameManager) {}
