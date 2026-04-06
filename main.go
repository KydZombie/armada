package main

import (
	"github.com/KydZombie/armada/core"
	"github.com/KydZombie/armada/game"
	"github.com/gen2brain/raylib-go/raylib"
)

func initWindow(config core.Config) {
	if config.VSync {
		rl.SetConfigFlags(rl.FlagVsyncHint)
	}

	rl.InitWindow(config.ScreenWidth, config.ScreenHeight, "Armada")
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())))
}

func main() {
	config := core.Config{
		ScreenWidth:  1280,
		ScreenHeight: 720,
		VSync:        false,

		Debug: true,
	}
	initWindow(config)
	defer rl.CloseWindow()

	gameManager := core.NewGameManager(config, &game.MainMenuScreen{})

	for !rl.WindowShouldClose() {
		gameManager.RunFrame()

		if gameManager.ShouldQuit {
			break
		}
	}

}
