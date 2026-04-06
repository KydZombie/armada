package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Screen interface {
	Update(gm *GameManager)
	Draw(gm *GameManager)
	DrawUI(gm *GameManager)
}

type MainMenuScreen struct{}

func (s *MainMenuScreen) Update(gm *GameManager) {
	if rl.IsKeyPressed(rl.KeySpace) {
		gm.SetScreen(NewGameScreen())
	}
}
func (s *MainMenuScreen) Draw(gm *GameManager) {
	rl.ClearBackground(rl.Black)
	rl.DrawText("Press space to start the game.", 64, 64, 64, rl.White)
}
func (s *MainMenuScreen) DrawUI(gm *GameManager) {}

type SettingsScreen struct {
	previousScreen Screen
}

func NewSettingsScreen(currentScreen Screen) *SettingsScreen {
	return &SettingsScreen{
		previousScreen: currentScreen,
	}
}

func (s *SettingsScreen) Update(gm *GameManager) {}
func (s *SettingsScreen) Draw(gm *GameManager)   {}
func (s *SettingsScreen) DrawUI(gm *GameManager) {}
