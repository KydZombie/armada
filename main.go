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
			Resizeable:   false,
			MasterVolume: 1.0,
			MusicVolume:  1.0,
			SFXVolume:    1.0,

			Debug: false,
		},
	)

	gameManager.CreateRaylibWindow()
	defer func() {
		for _, font := range gameManager.Fonts {
			rl.UnloadFont(font)
		}
		for _, tex := range gameManager.Textures {
			rl.UnloadTexture(tex)
		}
		rl.CloseWindow()
	}()

	gameManager.Fonts = map[string]rl.Font{
		"dh":   rl.LoadFont("assets/fonts/doublehomicide.ttf"),
		"ec_b": rl.LoadFont("assets/fonts/entercommand-b.ttf"),
		"ec_i": rl.LoadFont("assets/fonts/entercommand-i.ttf"),
		"ec":   rl.LoadFont("assets/fonts/entercommand.ttf"),
	}

	gameManager.Textures = map[string]rl.Texture2D{
		"ENG":      rl.LoadTexture("assets/textures/battle/system/engine.png"),
		"LIF":      rl.LoadTexture("assets/textures/battle/system/life.png"),
		"MED":      rl.LoadTexture("assets/textures/battle/system/medbay.png"),
		"PIL":      rl.LoadTexture("assets/textures/battle/system/pilot.png"),
		"SHD":      rl.LoadTexture("assets/textures/battle/system/shield.png"),
		"WPN":      rl.LoadTexture("assets/textures/battle/system/weapon.png"),
		"layout":   rl.LoadTexture("assets/textures/battle/layout.png"),
		"terminal": rl.LoadTexture("assets/textures/battle/terminal.png"),

		"back":       rl.LoadTexture("assets/textures/main/back.png"),
		"background": rl.LoadTexture("assets/textures/main/background.png"),
		"options":    rl.LoadTexture("assets/textures/main/options.png"),
		"title":      rl.LoadTexture("assets/textures/main/title.png"),
	}

	gameManager.SetScreen(game.NewMainMenuScreen())

	for !rl.WindowShouldClose() {
		gameManager.RunFrame()

		if gameManager.ShouldQuit {
			break
		}
	}
}
