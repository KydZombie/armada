package game

import (
	"github.com/KydZombie/armada/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type MainMenuScreen struct{}

func (s *MainMenuScreen) Update(gm *core.GameManager) {
	if rl.IsKeyPressed(rl.KeySpace) {
		gm.SetScreen(NewGameScreen())
	}
}
func (s *MainMenuScreen) Draw(gm *core.GameManager) {
	rl.ClearBackground(rl.Black)
	rl.DrawText("Press space to start the game.", 64, 64, 64, rl.White)
}
func (s *MainMenuScreen) DrawUI(gm *core.GameManager) {}

type SettingsScreen struct {
	previousScreen core.Screen
}

func NewSettingsScreen(currentScreen core.Screen) *SettingsScreen {
	return &SettingsScreen{
		previousScreen: currentScreen,
	}
}

func (s *SettingsScreen) Update(gm *core.GameManager) {}
func (s *SettingsScreen) Draw(gm *core.GameManager)   {}
func (s *SettingsScreen) DrawUI(gm *core.GameManager) {}
