package game

import (
	"github.com/KydZombie/armada/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type TutorialScreen struct {
	previousScreen core.Screen
}

func NewTutorialScreen(currentScreen core.Screen) *TutorialScreen {
	return &TutorialScreen{
		previousScreen: currentScreen,
	}
}

func (s *TutorialScreen) ResizeScreen(gm *core.GameManager) {}

func (s *TutorialScreen) UpdateScreen(gm *core.GameManager) {
	if rl.IsKeyPressed(rl.KeyEscape) || rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeyKpEnter) {
		gm.SetScreen(s.previousScreen)
		return
	}

	rect := rl.Rectangle{
		X:      float32(gm.NativeWidth)/2 - 175,
		Y:      float32(gm.NativeHeight) - 110,
		Width:  350,
		Height: 90,
	}
	if rl.CheckCollisionPointRec(gm.GetMouse(), rect) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		gm.SetScreen(s.previousScreen)
	}
}

func (s *TutorialScreen) DrawScreen(gm *core.GameManager) {
	rl.DrawTexture(gm.Textures["background"], 0, 0, rl.Color{R: 120, G: 120, B: 120, A: 255})

	rl.DrawTexture(gm.Textures["back"], 0, 0, rl.White)

	title := "Tutorial"

	sizeTitle := rl.MeasureTextEx(gm.Fonts["dh"], title, 100, 2)

	pos := rl.Vector2{
		X: float32(gm.NativeWidth)/2 - sizeTitle.X/2,
		Y: 70,
	}

	DrawEngravedText(
		gm.Fonts["dh"],
		title,
		pos,
		100,
		2,
		rl.White,
	)

	lines := []string{
		"1. Keep your ship running while under pressure.",
		"2. Enter commands in the terminal window during gameplay.",
		"3. Watch enemy and train windows to react quickly.",
		"4. Use command 'help' in-game to learn available actions.",
	}

	startY := int32(190)
	for i, line := range lines {
		posL := rl.Vector2{
			X: 120,
			Y: float32(startY + int32(i*48)),
		}

		DrawEngravedText(
			gm.Fonts["dh"],
			line,
			posL,
			30,
			2,
			rl.White,
		)
	}

	rect := rl.Rectangle{
		X:      float32(gm.NativeWidth)/2 - 175,
		Y:      float32(gm.NativeHeight) - 110,
		Width:  350,
		Height: 90,
	}

	hovered := rl.CheckCollisionPointRec(gm.GetMouse(), rect)
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

func (s *TutorialScreen) DrawScreenUI(gm *core.GameManager) {}
