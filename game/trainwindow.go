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
	mousePos := rl.GetMousePosition()
	if !rl.CheckCollisionPointRec(mousePos, t.GetBounds()) {
		return false
	}

	if !rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		return false
	}

	clickedCharacter := t.characterAtMouse(state, mousePos)
	if clickedCharacter != nil {
		state.SelectedCharacterIndex = clickedCharacter.Id
		return true
	}

	if state.SelectedCharacterIndex < 0 || state.SelectedCharacterIndex >= len(state.Train.Characters) {
		return true
	}

	target, ok := t.roomPosAtMouse(state, mousePos)
	if !ok {
		return true
	}

	character := state.Train.Characters[state.SelectedCharacterIndex]
	if _, err := state.Train.MoveCharacter(character, target); err != nil {
		gm.ErrLog.Printf("train movement failed: %v", err)
	}

	return true
}

func (t TrainWindow) UpdateWindow(gm *core.GameManager, state *Game) {

}

func (t TrainWindow) trainOffset() rl.Vector2 {
	bounds := t.GetBounds()
	return rl.Vector2{X: bounds.X + 16.0, Y: bounds.Y + 48.0}
}

func (t TrainWindow) tileSize() float32 {
	bounds := t.GetBounds()

	// Fit the full train layout (about 22-23 tiles wide) inside current window width.
	widthBudget := bounds.Width - 56
	tileByWidth := widthBudget / 23.0

	// Keep room bars and the bottom stats strip visible as well.
	heightBudget := bounds.Height - 156
	tileByHeight := heightBudget / 3.0

	tile := float32(48.0)
	if tileByWidth < tile {
		tile = tileByWidth
	}
	if tileByHeight < tile {
		tile = tileByHeight
	}
	if tile < 18 {
		tile = 18
	}

	return tile
}

func (t TrainWindow) roomBounds(room Room) rl.Rectangle {
	trainOffset := t.trainOffset()
	tileSize := t.tileSize()
	return rl.Rectangle{
		X:      trainOffset.X + room.Pos.X*tileSize,
		Y:      trainOffset.Y + room.Pos.Y*tileSize,
		Width:  float32(room.Width) * tileSize,
		Height: float32(room.Height) * tileSize,
	}
}

func (t TrainWindow) hallwayBounds(leftRoom Room, rightRoom Room) (rl.Rectangle, bool) {
	leftBounds := t.roomBounds(leftRoom)
	rightBounds := t.roomBounds(rightRoom)

	hallwayX := leftBounds.X + leftBounds.Width
	hallwayWidth := rightBounds.X - hallwayX
	if hallwayWidth <= 0 {
		return rl.Rectangle{}, false
	}

	tileSize := t.tileSize()
	return rl.Rectangle{
		X:      hallwayX,
		Y:      leftBounds.Y + tileSize,
		Width:  hallwayWidth,
		Height: tileSize,
	}, true
}

func (t TrainWindow) roomPosAtMouse(state *Game, mousePos rl.Vector2) (RoomPos, bool) {
	for _, room := range state.Train.Rooms {
		bounds := t.roomBounds(room)
		if !rl.CheckCollisionPointRec(mousePos, bounds) {
			continue
		}

		tileSize := t.tileSize()
		relativeX := int((mousePos.X - bounds.X) / tileSize)
		relativeY := int((mousePos.Y - bounds.Y) / tileSize)
		if room.HasTile(relativeX, relativeY) {
			return RoomPos{RoomId: room.Id, X: relativeX, Y: relativeY}, true
		}
	}

	return RoomPos{}, false
}

func (t TrainWindow) characterAtMouse(state *Game, mousePos rl.Vector2) *Character {
	for _, character := range state.Train.Characters {
		characterPos := t.characterWorldPosition(state, character)
		characterRect := rl.Rectangle{
			X:      characterPos.X,
			Y:      characterPos.Y,
			Width:  t.tileSize(),
			Height: t.tileSize(),
		}
		if rl.CheckCollisionPointRec(mousePos, characterRect) {
			return character
		}
	}

	return nil
}

func (t TrainWindow) roomTileWorldPosition(room Room, tileX, tileY int) rl.Vector2 {
	tileSize := t.tileSize()
	offset := t.trainOffset()
	return rl.Vector2{
		X: offset.X + (room.Pos.X+float32(tileX))*tileSize,
		Y: offset.Y + (room.Pos.Y+float32(tileY))*tileSize,
	}
}

func (t TrainWindow) characterWorldPosition(state *Game, character *Character) rl.Vector2 {
	if !character.IsMoving || len(character.MovementPath) == 0 {
		room, _ := state.Train.GetRoom(character.Pos.RoomId)
		return t.roomTileWorldPosition(*room, character.Pos.X, character.Pos.Y)
	}

	if character.CurrentPathIndex >= len(character.MovementPath)-1 {
		finalPos := character.MovementPath[len(character.MovementPath)-1]
		room, _ := state.Train.GetRoom(finalPos.RoomId)
		return t.roomTileWorldPosition(*room, finalPos.X, finalPos.Y)
	}

	currentPos := character.MovementPath[character.CurrentPathIndex]
	nextPos := character.MovementPath[character.CurrentPathIndex+1]
	currentRoom, _ := state.Train.GetRoom(currentPos.RoomId)
	nextRoom, _ := state.Train.GetRoom(nextPos.RoomId)

	currentWorld := t.roomTileWorldPosition(*currentRoom, currentPos.X, currentPos.Y)
	nextWorld := t.roomTileWorldPosition(*nextRoom, nextPos.X, nextPos.Y)

	if currentPos.RoomId != nextPos.RoomId {
		return rl.Vector2{
			X: currentWorld.X + (nextWorld.X-currentWorld.X)*character.AnimationProgress,
			Y: currentWorld.Y + (nextWorld.Y-currentWorld.Y)*character.AnimationProgress,
		}
	}

	moveX := nextWorld.X - currentWorld.X
	moveY := nextWorld.Y - currentWorld.Y

	animated := currentWorld
	if moveX != 0 {
		if character.AnimationProgress < 0.5 {
			animated.X = currentWorld.X + moveX*(character.AnimationProgress*2.0)
		} else {
			animated.X = nextWorld.X
			animated.Y = currentWorld.Y + moveY*((character.AnimationProgress-0.5)*2.0)
		}
	} else if moveY != 0 {
		animated.Y = currentWorld.Y + moveY*character.AnimationProgress
	}

	return animated
}

func (t TrainWindow) DrawWindow(gm *core.GameManager, state *Game) {
	bounds := t.GetBounds()
	rl.DrawRectangleRec(bounds, rl.Blue)
	rl.BeginScissorMode(int32(bounds.X), int32(bounds.Y), int32(bounds.Width), int32(bounds.Height))

	// TODO: Use sprites for train rendering

	tileSize := t.tileSize()
	const roomBorderThickness float32 = 3.0
	const roomLabelFontSize int32 = 28
	const roomBarTextSize int32 = 16
	const roomBarHeight float32 = 18
	const roomBarSpacing float32 = 6

	for roomIdx := 0; roomIdx < len(state.Train.Rooms)-1; roomIdx++ {
		hallwayBounds, ok := t.hallwayBounds(state.Train.Rooms[roomIdx], state.Train.Rooms[roomIdx+1])
		if !ok {
			continue
		}

		hallwayCols := int(hallwayBounds.Width/tileSize + 0.5)
		hallwayRows := int(hallwayBounds.Height/tileSize + 0.5)
		if hallwayCols < 1 {
			hallwayCols = 1
		}
		if hallwayRows < 1 {
			hallwayRows = 1
		}
		for x := 0; x < hallwayCols; x++ {
			for y := 0; y < hallwayRows; y++ {
				tileBounds := rl.Rectangle{
					X:      hallwayBounds.X + float32(x)*tileSize,
					Y:      hallwayBounds.Y + float32(y)*tileSize,
					Width:  tileSize,
					Height: tileSize,
				}
				rl.DrawRectangleRec(tileBounds, rl.NewColor(215, 215, 205, 255))
				rl.DrawRectangleLinesEx(tileBounds, 1.5, rl.DarkGray)
			}
		}
		rl.DrawRectangleLinesEx(hallwayBounds, 2.0, rl.Black)
	}

	for _, room := range state.Train.Rooms {
		roomBounds := t.roomBounds(room)
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
		spaceAroundSideOfDoor := tileSize / 6
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
		rl.DrawText(room.System.ShortName(), int32(roomBounds.X)+34, int32(labelY)+5, 16, rl.DarkBlue)

		barWidth := roomBounds.Width - 8
		barX := roomBounds.X + 4
		healthBarY := roomBounds.Y + roomBounds.Height + 6
		damageBarY := healthBarY + roomBarHeight + roomBarSpacing
		maxDamageDisplay := room.AttackPower
		if maxDamageDisplay < 5 {
			maxDamageDisplay = 5
		}

		barLabel := "ATK"
		barValue := room.AttackPower
		if room.System.Type != SystemWeapons {
			barLabel = "SYS"
			barValue = 0
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
			barLabel,
			barValue,
			maxDamageDisplay,
			rl.Gold,
			roomBarTextSize,
		)
	}

	for _, character := range state.Train.Characters {
		characterPos := t.characterWorldPosition(state, character)

		var renderColor rl.Color
		if state.SelectedCharacterIndex == character.Id {
			renderColor = rl.Green
		} else {
			renderColor = rl.DarkGray
		}
		rl.DrawCircleV(rl.Vector2AddValue(characterPos, tileSize/2), tileSize/3, renderColor)

		const fontSize int32 = 18
		text := fmt.Sprint(character.Id + 1)
		textWidth := rl.MeasureText(text, fontSize)

		rl.DrawText(
			text,
			int32(characterPos.X+tileSize/2-float32(textWidth)/2),
			int32(characterPos.Y+tileSize/2-float32(fontSize)/2),
			fontSize,
			rl.White,
		)
	}

	statsHeight := float32(58)
	statsBounds := rl.Rectangle{
		X:      bounds.X + 8,
		Y:      bounds.Y + bounds.Height - statsHeight - 8,
		Width:  bounds.Width - 16,
		Height: statsHeight,
	}

	rl.DrawRectangleRec(statsBounds, rl.Fade(rl.Black, 0.48))
	rl.DrawRectangleLinesEx(statsBounds, 2, rl.Fade(rl.White, 0.35))

	hullText := fmt.Sprintf("Hull %d/%d", state.Train.Health, state.Train.MaxHealth)
	defenseText := fmt.Sprintf("Shields %d   Evasion %d%%   Weapons Ready %d/%d", state.Train.ShieldLayers(), state.Train.EvasionChance(), state.Train.ReadyWeapons(), len(state.Train.Weapons))
	medbayText := fmt.Sprintf("Medbay +%d/tick", state.Train.MedbayHealingPerTick())
	lifeSupportText := "Life Support online"
	if !state.Train.LifeSupportOperational() {
		lifeSupportText = fmt.Sprintf("Life Support offline (%d/tick)", state.Train.LifeSupportDamagePerTick())
	}

	rl.DrawText(hullText, int32(statsBounds.X+10), int32(statsBounds.Y+6), 22, rl.White)
	rl.DrawText(defenseText, int32(statsBounds.X+200), int32(statsBounds.Y+9), 18, rl.LightGray)
	rl.DrawText(medbayText, int32(statsBounds.X+10), int32(statsBounds.Y+32), 16, rl.Green)
	rl.DrawText(lifeSupportText, int32(statsBounds.X+170), int32(statsBounds.Y+32), 16, rl.Orange)

	rl.EndScissorMode()
}

func (t TrainWindow) DrawWindowUI(gm *core.GameManager, state *Game) {
	// Instructions removed per request
}
