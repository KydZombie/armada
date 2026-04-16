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
	const roomBorderThickness float32 = 3.0

	for _, room := range state.Train.Rooms {
		roomBounds := rl.Rectangle{
			X:      trainOffset.X + room.Pos.X*tileSize,
			Y:      trainOffset.Y + room.Pos.Y*tileSize,
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

		rl.DrawRectangleLinesEx(roomBounds, roomBorderThickness, rl.Black)
		doorThickness := tileSize / 8
		spaceAroundSideOfDoor := tileSize / 6 // Essentially inverse of door width
		for _, door := range room.Doors {
			var doorBounds rl.Rectangle
			switch door.Facing {
			case core.FacingLeft:
				doorBounds = rl.Rectangle{
					X:      roomBounds.X + float32(door.X)*tileSize + roomBorderThickness,
					Y:      roomBounds.Y + float32(door.Y)*tileSize + spaceAroundSideOfDoor,
					Width:  doorThickness,
					Height: tileSize - 2*spaceAroundSideOfDoor,
				}
			case core.FacingRight:
				doorBounds = rl.Rectangle{
					X:      roomBounds.X + float32(door.X+1)*tileSize - doorThickness - roomBorderThickness,
					Y:      roomBounds.Y + float32(door.Y)*tileSize + spaceAroundSideOfDoor,
					Width:  doorThickness,
					Height: tileSize - 2*spaceAroundSideOfDoor,
				}
			default:
				gm.ErrLog.Println("core.FacingUp and core.FacingDown door rendering is not implemented yet.")
				doorBounds = rl.Rectangle{X: 0, Y: 0, Width: 0, Height: 0}
			}

			rl.DrawRectangleRec(doorBounds, rl.Orange)
		}

		rl.DrawText(string([]rune{room.GetRune()}), int32(roomBounds.X)+4, int32(roomBounds.Y)+4, 24, rl.Black)
	}
}

func (t TrainWindow) DrawWindowUI(gm *core.GameManager, state *Game) {
}
