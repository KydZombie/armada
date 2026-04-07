package core

import (
	"fmt"
	"log"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameManager struct {
	Config

	Log    *log.Logger
	ErrLog *log.Logger

	// Screen should never be nil
	Screen    Screen
	DeltaTime float32

	ShouldQuit bool
}

func NewGameManager(config Config, startScreen Screen) *GameManager {
	return &GameManager{
		Config: config,

		Log:    log.New(os.Stdout, "", log.LstdFlags),
		ErrLog: log.New(os.Stderr, "", log.LstdFlags),

		Screen:    startScreen,
		DeltaTime: 0,
	}
}

func (gm *GameManager) Quit() {
	gm.ShouldQuit = true
}

func (gm *GameManager) SetScreen(screen Screen) {
	gm.Screen = screen
}

func (gm *GameManager) RunFrame() {
	gm.DeltaTime = rl.GetFrameTime()

	if rl.IsKeyPressed(rl.KeyF3) {
		gm.Debug = !gm.Debug
	}

	gm.Screen.UpdateScreen(gm)

	rl.BeginDrawing()
	defer rl.EndDrawing()

	gm.Screen.DrawScreen(gm)
	gm.Screen.DrawScreenUI(gm)
	if gm.Debug {
		dtText := fmt.Sprintf("FrameTime: %.4f", gm.DeltaTime)
		rl.DrawRectangle(0, 0, 190, 40, rl.Black)
		rl.DrawFPS(4, 4)
		rl.DrawText(dtText, 4, 20, 20, rl.DarkGreen)
	}
}
