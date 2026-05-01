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
			NativeWidth:  1280,
			NativeHeight: 720,
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

	gameManager.Shader = rl.LoadShader("assets/shaders/scan.vert", "assets/shaders/scan.frag")

	gameManager.Fonts = map[string]rl.Font{
		"dh":   rl.LoadFontEx("assets/fonts/doublehomicide.ttf", 64, nil),
		"ec-b": rl.LoadFontEx("assets/fonts/entercommand-b.ttf", 64, nil),
		"ec-i": rl.LoadFontEx("assets/fonts/entercommand-i.ttf", 64, nil),
		"ec":   rl.LoadFontEx("assets/fonts/entercommand.ttf", 64, nil),
	}

	gameManager.Textures = map[string]rl.Texture2D{
		"ENG":      rl.LoadTexture("assets/textures/battle/system/engine.png"),
		"LIF":      rl.LoadTexture("assets/textures/battle/system/life.png"),
		"MED":      rl.LoadTexture("assets/textures/battle/system/medbay.png"),
		"PIL":      rl.LoadTexture("assets/textures/battle/system/pilot.png"),
		"SHD":      rl.LoadTexture("assets/textures/battle/system/shield.png"),
		"WPN":      rl.LoadTexture("assets/textures/battle/system/weapon.png"),
		"enemyB":   rl.LoadTexture("assets/textures/battle/enemyB.png"),
		"enemyF":   rl.LoadTexture("assets/textures/battle/enemyF.png"),
		"layout":   rl.LoadTexture("assets/textures/battle/layout.png"),
		"terminal": rl.LoadTexture("assets/textures/battle/terminal.png"),
		"trainS":   rl.LoadTexture("assets/textures/battle/trainS.png"),
		"trainT":   rl.LoadTexture("assets/textures/battle/trainT.png"),

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
