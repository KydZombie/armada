package game

import (
	"fmt"

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
	const roomLabelFontSize int32 = 24
	const roomBarTextSize int32 = 14
	const roomBarHeight float32 = 18
	const roomBarSpacing float32 = 6

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

		for _, character := range room.GetCharacters(state.Train) {
			pos := rl.Vector2{
				X: roomBounds.X + float32(character.Pos.X)*tileSize,
				Y: roomBounds.Y + float32(character.Pos.Y)*tileSize,
			}

			var renderColor rl.Color
			if state.SelectedCharacterIndex == character.Id {
				renderColor = rl.Green
			} else {
				renderColor = rl.DarkGray
			}
			rl.DrawCircleV(rl.Vector2AddValue(pos, tileSize/2), tileSize/3, renderColor)

			const fontSize int32 = 16
			text := fmt.Sprint(character.Id + 1)
			textWidth := rl.MeasureText(text, fontSize)

			rl.DrawText(
				text,
				int32(pos.X+tileSize/2-float32(textWidth)/2),
				int32(pos.Y+tileSize/2-float32(fontSize)/2),
				fontSize,
				rl.White,
			)
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

		labelY := roomBounds.Y + 4
		rl.DrawText(string([]rune{room.GetRune()}), int32(roomBounds.X)+4, int32(labelY), roomLabelFontSize, rl.Black)
		rl.DrawText(room.System.ShortName(), int32(roomBounds.X)+28, int32(labelY)+4, 14, rl.DarkBlue)

		barWidth := roomBounds.Width - 8
		barX := roomBounds.X + 4
		healthBarY := roomBounds.Y + roomBounds.Height + 6
		damageBarY := healthBarY + roomBarHeight + roomBarSpacing
		maxDamageDisplay := room.AttackPower
		if maxDamageDisplay < 5 {
			maxDamageDisplay = 5
		}

		drawStatBar(
			rl.Rectangle{
				X:      barX,
				Y:      healthBarY,
				Width:  barWidth,
				Height: roomBarHeight,
			},
			"HP",
			room.Health,
			room.MaxHealth,
			rl.Red,
			roomBarTextSize,
		)

		drawStatBar(
			rl.Rectangle{
				X:      barX,
				Y:      damageBarY,
				Width:  barWidth,
				Height: roomBarHeight,
			},
			"DMG",
			room.AttackPower,
			maxDamageDisplay,
			rl.Gold,
			roomBarTextSize,
		)
	}
}

func (t TrainWindow) DrawWindowUI(gm *core.GameManager, state *Game) {
}
