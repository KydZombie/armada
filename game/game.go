package game

import (
	"github.com/KydZombie/armada/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	Train *Train
	Enemy Enemy

	PlayerWeapon Weapon

	RoomEnemies            []Enemy
	Enemies                []Enemy
	SelectedRoom           int
	SelectionPopupText     string
	SelectionPopupFrames   int
	enemyTimerFrameCounter int

	windows []core.Window[Game]
}

func NewGameScreen(gm *core.GameManager) *Game {
	train := NewTrain(100)
	roomEnemies := []Enemy{
		NewBasicEnemy("Steel Crawler", 20, 3),
		NewBasicEnemy("Iron crawler", 14, 2),
		nil,
	}

	gs := &Game{
		Train:        train,
		PlayerWeapon: NewWeapon("Cannon", 1, 1),
		RoomEnemies:  roomEnemies,
		Enemies:      roomEnemies,
		SelectedRoom: 0,

		windows: []core.Window[Game]{},
	}

	gs.syncSelectedRoom()

	const windowMargin = 16.0

	terminal := NewTerminalWindow(
		func(gm *core.GameManager) rl.Rectangle {
			return rl.Rectangle{
				X:      windowMargin,
				Y:      float32(gm.ScreenHeight)/2.0 + windowMargin,
				Width:  float32(gm.ScreenWidth)/2.0 - windowMargin,
				Height: float32(gm.ScreenHeight)/2.0 - windowMargin*2,
			}
		},
		gm,
		initializeCommands(),
	)

	gs.windows = append(gs.windows, terminal)

	gs.windows = append(gs.windows, NewTrainWindow(
		func(gm *core.GameManager) rl.Rectangle {
			return rl.Rectangle{
				X:      windowMargin,
				Y:      windowMargin,
				Width:  float32(gm.ScreenWidth) - windowMargin*2,
				Height: float32(gm.ScreenHeight)/2.0 - windowMargin,
			}
		},
		gm,
	))

	gs.windows = append(gs.windows, NewBattleWindow(
		func(gm *core.GameManager) rl.Rectangle {
			return rl.Rectangle{
				X:      float32(gm.ScreenWidth)/2.0 + windowMargin,
				Y:      float32(gm.ScreenHeight)/2.0 + windowMargin,
				Width:  float32(gm.ScreenWidth)/2.0 - windowMargin*2,
				Height: float32(gm.ScreenHeight)/2.0 - windowMargin*2,
			}
		},
		gm,
	))

	return gs
}

func (g *Game) ResizeScreen(gm *core.GameManager) {
	g.UpdateWindowSizes(gm)
}

func (g *Game) UpdateScreen(gm *core.GameManager) {
	inputCaptured := false

	if g.SelectionPopupFrames > 0 {
		g.SelectionPopupFrames--
	}

	for _, window := range g.windows {
		// If a window captures the input, other windows should not read any input
		if !inputCaptured && window.HandleInput(gm, g) {
			inputCaptured = true
		}

		window.UpdateWindow(gm, g)
	}
}

func (g *Game) DrawScreen(gm *core.GameManager) {
	rl.ClearBackground(rl.DarkBlue)

	for _, window := range g.windows {
		window.DrawWindow(gm, g)
	}
}

func (g *Game) DrawScreenUI(gm *core.GameManager) {
	//var buttonText string
	//if g.moving {
	//	buttonText = "Stop moving"
	//} else {
	//	buttonText = "Start moving"
	//}
	//
	//if rg.Button(rl.Rectangle{
	//	X:      0,
	//	Y:      50,
	//	Width:  100,
	//	Height: 60,
	//}, buttonText) {
	//	g.moving = !g.moving
	//}

	for _, window := range g.windows {
		window.DrawWindowUI(gm, g)
	}
}

func (g *Game) UpdateWindowSizes(gm *core.GameManager) {
	for _, window := range g.windows {
		window.UpdateWindowSize(gm)
	}
}

// syncSelectedRoom keeps the old Enemy field aligned with the currently
// selected room so older single-enemy code can keep working.
func (g *Game) syncSelectedRoom() {
	if len(g.RoomEnemies) == 0 {
		g.Enemy = nil
		g.SelectedRoom = 0
		return
	}

	if g.SelectedRoom < 0 {
		g.SelectedRoom = 0
	}

	if g.SelectedRoom >= len(g.RoomEnemies) {
		g.SelectedRoom = len(g.RoomEnemies) - 1
	}

	selectedEnemy := g.RoomEnemies[g.SelectedRoom]
	g.Enemy = selectedEnemy
}

func (g *Game) SelectRoom(roomIndex int) {
	g.SelectedRoom = roomIndex
	g.syncSelectedRoom()
}
