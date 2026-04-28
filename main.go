package main

import (
	"github.com/KydZombie/armada/core"
	"github.com/KydZombie/armada/game"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	gameManager := core.NewGameManager(
		"Armada",
		core.Config{
			ScreenWidth:  1280,
			ScreenHeight: 720,
			VSync:        false,
			Resizeable:   true,
			MasterVolume: 1.0,
			MusicVolume:  1.0,
			SFXVolume:    1.0,

			Debug: false,
		},
	)

	gameManager.CreateRaylibWindow()
	defer rl.CloseWindow()

	gameManager.SetScreen(game.NewMainMenuScreen())

	for !rl.WindowShouldClose() {
		gameManager.RunFrame()

		if gameManager.ShouldQuit {
			break
		}
	}
}
