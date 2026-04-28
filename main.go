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

	gameManager.Fonts = map[string]rl.Font{
		"dh":   rl.LoadFont("assets/fonts/doublehomicide.ttf"),
		"ec_b": rl.LoadFont("assets/fonts/entercommand-b.ttf"),
		"ec_i": rl.LoadFont("assets/fonts/entercommand-i.ttf"),
		"ec":   rl.LoadFont("assets/fonts/entercommand.ttf"),
	}

	gameManager.Textures = map[string]rl.Texture2D{
		"engine": rl.LoadTexture("assets/menus/battle/system/engine.png"),
		"life":   rl.LoadTexture("assets/menus/battle/system/life.png"),
		"medbay": rl.LoadTexture("assets/menus/battle/system/medbay.png"),
		"pilot":  rl.LoadTexture("assets/menus/battle/system/pilot.png"),
		"shield": rl.LoadTexture("assets/menus/battle/system/shield.png"),
		"weapon": rl.LoadTexture("assets/menus/battle/system/weapon.png"),

		"optionsM": rl.LoadTexture("assets/menus/main/optionsM.png"),
		"titleM":   rl.LoadTexture("assets/menus/main/titleM.png"),

		"back":       rl.LoadTexture("assets/misc/back.png"),
		"background": rl.LoadTexture("assets/misc/background.png"),
	}

	gameManager.SetScreen(game.NewMainMenuScreen())

	for !rl.WindowShouldClose() {
		gameManager.RunFrame()

		if gameManager.ShouldQuit {
			break
		}
	}
}
