package core

import (
	"fmt"
	"log"
	"math"
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

	Shader        rl.Shader
	RenderTexture rl.RenderTexture2D

	Fonts    map[string]rl.Font
	Textures map[string]rl.Texture2D
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
	gm.RenderTexture = rl.LoadRenderTexture(gm.Config.ScreenWidth, gm.Config.ScreenHeight)
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())))

	if gm.Config.Fullscreen {
		rl.ToggleBorderlessWindowed()
		gm.ScreenWidth = int32(rl.GetScreenWidth())
		gm.ScreenHeight = int32(rl.GetScreenHeight())
	}
}

func (gm *GameManager) SetVSync(enabled bool) {
	gm.VSync = enabled
	if rl.IsWindowFullscreen() {
		return
	}
	if enabled {
		rl.SetWindowState(rl.FlagVsyncHint)
	} else {
		rl.ClearWindowState(rl.FlagVsyncHint)
	}
}

func (gm *GameManager) SetFullscreen(enabled bool) {
	if rl.IsWindowState(rl.FlagBorderlessWindowedMode) == enabled {
		gm.Fullscreen = enabled
		return
	}

	rl.ToggleBorderlessWindowed()
	gm.Fullscreen = enabled
	gm.ScreenWidth = int32(rl.GetScreenWidth())
	gm.ScreenHeight = int32(rl.GetScreenHeight())

	if gm.Screen != nil {
		gm.Screen.ResizeScreen(gm)
	}
}

func (gm *GameManager) SetScreen(screen Screen) {
	gm.Screen = screen
}

func (gm *GameManager) GetMouse() rl.Vector2 {
	screenW := float32(rl.GetScreenWidth())
	screenH := float32(rl.GetScreenHeight())

	scaleX := screenW / float32(gm.NativeWidth)
	scaleY := screenH / float32(gm.NativeHeight)

	scale := float32(math.Min(float64(scaleX), float64(scaleY)))

	destW := float32(gm.NativeWidth) * scale
	//destH := float32(gm.NativeHeight) * scale

	offsetX := (screenW - destW) / 2
	offsetY := float32(0) //(screenH - destH) / 2

	mouse := rl.GetMousePosition()

	return rl.Vector2{
		X: (mouse.X - offsetX) / scale,
		Y: (mouse.Y - offsetY) / scale,
	}
}

func (gm *GameManager) RunFrame() {
	gm.DeltaTime = rl.GetFrameTime()

	if rl.IsKeyPressed(rl.KeyF3) {
		gm.Debug = !gm.Debug
	}

	if rl.IsWindowResized() {
		gm.ScreenWidth = int32(rl.GetScreenWidth())
		gm.ScreenHeight = int32(rl.GetScreenHeight())

		rl.UnloadRenderTexture(gm.RenderTexture)
		gm.RenderTexture = rl.LoadRenderTexture(gm.ScreenWidth, gm.ScreenHeight)

		gm.Screen.ResizeScreen(gm)
	}

	gm.Screen.UpdateScreen(gm)

	rl.BeginTextureMode(gm.RenderTexture)
	rl.ClearBackground(rl.Black)

	gm.Screen.DrawScreen(gm)
	gm.Screen.DrawScreenUI(gm)
	if gm.Debug {
		dtText := fmt.Sprintf("FrameTime: %.4f", gm.DeltaTime)
		rl.DrawRectangle(0, 0, 190, 40, rl.Black)
		rl.DrawFPS(4, 4)
		rl.DrawText(dtText, 4, 20, 20, rl.DarkGreen)
	}

	rl.EndTextureMode()

	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	screenW := float32(rl.GetScreenWidth())
	screenH := float32(rl.GetScreenHeight())

	scaleX := screenW / float32(gm.NativeWidth)
	scaleY := screenH / float32(gm.NativeHeight)

	scale := float32(math.Min(float64(scaleX), float64(scaleY)))

	destW := float32(gm.NativeWidth) * scale
	destH := float32(gm.NativeHeight) * scale

	offsetX := (screenW - destW) / 2
	offsetY := float32(0) //(screenH - destH) / 2

	rl.BeginShaderMode(gm.Shader)

	rl.DrawTexturePro(
		gm.RenderTexture.Texture,
		rl.NewRectangle(
			0,
			0,
			float32(gm.RenderTexture.Texture.Width),
			-float32(gm.RenderTexture.Texture.Height),
		),
		rl.NewRectangle(
			offsetX,
			offsetY,
			destW,
			destH,
		),
		rl.NewVector2(0, 0),
		0,
		rl.White,
	)

	rl.EndShaderMode()

	rl.EndDrawing()
}
