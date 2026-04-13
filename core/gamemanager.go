package core

import (
	"fmt"
	"log"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameManager struct {
	Config

	WindowTitle string

	Log    *log.Logger
	ErrLog *log.Logger

	// Screen should never be nil if the game is currently running
	Screen    Screen
	DeltaTime float32

	ShouldQuit bool
}

func NewGameManager(windowTitle string, config Config) *GameManager {
	return &GameManager{
		Config: config,

		Log:    log.New(os.Stdout, "", log.LstdFlags),
		ErrLog: log.New(os.Stderr, "", log.LstdFlags),

		DeltaTime: 0,
	}
}

func (gm *GameManager) Quit() {
	gm.ShouldQuit = true
}

func (gm *GameManager) CreateRaylibWindow() {
	if gm.Config.VSync {
		rl.SetConfigFlags(rl.FlagVsyncHint)
	}
	if gm.Config.Resizeable {
		rl.SetConfigFlags(rl.FlagWindowResizable)
	}

	rl.InitWindow(gm.Config.ScreenWidth, gm.Config.ScreenHeight, "Armada")
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())))
}

func (gm *GameManager) SetScreen(screen Screen) {
	gm.Screen = screen
}

func (gm *GameManager) RunFrame() {
	gm.DeltaTime = rl.GetFrameTime()

	if rl.IsKeyPressed(rl.KeyF3) {
		gm.Debug = !gm.Debug
	}

	if rl.IsWindowResized() {
		gm.ScreenWidth = int32(rl.GetScreenWidth())
		gm.ScreenHeight = int32(rl.GetScreenHeight())
		gm.Screen.ResizeScreen(gm)
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
