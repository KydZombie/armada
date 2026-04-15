package main

import (
	"github.com/KydZombie/armada/core"
	"github.com/KydZombie/armada/game"
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	gameManager := core.NewGameManager(
		"Armada",
		core.Config{
			ScreenWidth:  1280,
			ScreenHeight: 720,
			VSync:        false,
			Resizeable:   true,

			Debug: true,
		},
	)

	gameManager.CreateRaylibWindow()
	defer rl.CloseWindow()

	gameManager.SetScreen(&game.MainMenuScreen{})

	for !rl.WindowShouldClose() {
		gameManager.RunFrame()

		if gameManager.ShouldQuit {
			break
		}
	}
}
