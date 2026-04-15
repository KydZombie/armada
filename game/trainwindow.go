package game

import (
	"github.com/KydZombie/armada/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type TrainWindow struct {
	core.BaseWindow[Game]
}

func NewTrainWindow(sizeFunc func(gm *core.GameManager) rl.Rectangle, gm *core.GameManager) *TrainWindow {
	return &TrainWindow{
		BaseWindow: core.NewBaseWindow[Game](sizeFunc, gm, true),
	}
}

func (t TrainWindow) HandleInput(gm *core.GameManager, state *Game) bool {
	return false
}

func (t TrainWindow) UpdateWindow(gm *core.GameManager, state *Game) {

}

func (t TrainWindow) DrawWindow(gm *core.GameManager, state *Game) {
	rl.DrawRectangleRec(t.GetBounds(), rl.Blue)

	// TODO: Use sprites for train rendering

	bounds := t.GetBounds()
	trainOffset := rl.Vector2{
		X: bounds.X + 16.0,
		Y: bounds.Y + 48.0,
	}

	const tileSize float32 = 48.0

	for _, room := range state.Train.Rooms {
		roomBounds := rl.Rectangle{
			X:      trainOffset.X + float32(room.X)*tileSize,
			Y:      trainOffset.Y + float32(room.Y)*tileSize,
			Width:  float32(room.Width) * tileSize,
			Height: float32(room.Height) * tileSize,
		}
		for x := range room.Width {
			for y := range room.Height {
				tileBounds := rl.Rectangle{
					X:      roomBounds.X + float32(x)*tileSize,
					Y:      roomBounds.Y + float32(y)*tileSize,
					Width:  tileSize,
					Height: tileSize,
				}
				rl.DrawRectangleRec(tileBounds, rl.RayWhite)
				rl.DrawRectangleLinesEx(tileBounds, 2.0, rl.Gray)
			}
		}

		rl.DrawRectangleLinesEx(roomBounds, 3.0, rl.Black)
		for _, door := range room.Doors {
			doorBounds := rl.Rectangle{
				X:      roomBounds.X + float32(door.X)*tileSize + (tileSize / 4),
				Y:      roomBounds.Y + float32(door.Y)*tileSize + (tileSize / 4),
				Width:  tileSize / 2,
				Height: tileSize / 2,
			}
			rl.DrawRectangleRec(doorBounds, rl.Green)
		}
	}
}

func (t TrainWindow) DrawWindowUI(gm *core.GameManager, state *Game) {
}
