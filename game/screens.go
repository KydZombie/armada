package game

import (
	"github.com/KydZombie/armada/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type MainMenuScreen struct{}

func (s *MainMenuScreen) UpdateScreen(gm *core.GameManager) {
	if rl.IsKeyPressed(rl.KeySpace) {
		gm.SetScreen(NewGameScreen(gm))
	}
}

func (s *MainMenuScreen) DrawScreen(gm *core.GameManager) {
	rl.ClearBackground(rl.Black)
	rl.DrawText("Press space to start the game.", 64, 64, 64, rl.White)
}

func (s *MainMenuScreen) DrawScreenUI(gm *core.GameManager) {}

type SettingsScreen struct {
	previousScreen core.Screen
}

func NewSettingsScreen(currentScreen core.Screen) *SettingsScreen {
	return &SettingsScreen{
		previousScreen: currentScreen,
	}
}

func (s *SettingsScreen) UpdateScreen(gm *core.GameManager) {}
func (s *SettingsScreen) DrawScreen(gm *core.GameManager)   {}
func (s *SettingsScreen) DrawScreenUI(gm *core.GameManager) {}
