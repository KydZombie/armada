package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func initWindow(config Config) {
	if config.VSync {
		rl.SetConfigFlags(rl.FlagVsyncHint)
	}

	rl.InitWindow(config.ScreenWidth, config.ScreenHeight, "Armada")
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())))
}

func main() {
	config := Config{
		ScreenWidth:  1280,
		ScreenHeight: 720,
		VSync:        false,

		Debug: true,
	}
	initWindow(config)
	defer rl.CloseWindow()

	gameManager := NewGameManager(config, &MainMenuScreen{})

	for !rl.WindowShouldClose() {
		gameManager.runLoop()

		if gameManager.shouldQuit {
			break
		}
	}

}
