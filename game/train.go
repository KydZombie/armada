package game

import (
	"fmt"
	"strings"

	"github.com/KydZombie/armada/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Room struct {
	Id int

	Pos           rl.Vector2
	Width, Height int

	Doors []Door

	System System
}

func (room *Room) GetRune() rune {
	return rune('A' + room.Id)
}

func (room *Room) HasTile(x, y int) bool {
	return x >= 0 && x < room.Width && y >= 0 && y < room.Height
}

func (room *Room) GetDoorTo(roomId int) *Door {
	for i := range room.Doors {
		if room.Doors[i].GoesToRoom.RoomId == roomId {
			return &room.Doors[i]
		}
	}

	return nil
}

func (room *Room) GetCharacters(train *Train) []*Character {
	var characters []*Character
	for _, character := range train.Characters {
		if character.Pos.RoomId == room.Id {
			characters = append(characters, character)
		}
	}
	return characters
}

type RoomPos struct {
	RoomId int
	X, Y   int
}

type Door struct {
	X, Y   int
	Facing core.Facing

	GoesToRoom RoomPos
}

type Train struct {
	Health, MaxHealth int
	Rooms             []Room
	Characters        []*Character
}

func NewTrain(health int) *Train {
	train := &Train{
		Health:     health,
		MaxHealth:  health,
		Rooms:      make([]Room, 0),
		Characters: make([]*Character, 0),
	}

	train.addRoom(Room{
		Pos: rl.Vector2{
			X: 0,
			Y: 0,
		},

		Width:  4,
		Height: 3,
		Doors: []Door{
			{
				X:      3,
				Y:      1,
				Facing: core.FacingRight,
				GoesToRoom: RoomPos{
					RoomId: 1,
					X:      0,
					Y:      1,
				},
			},
		},
	})

	train.addRoom(Room{
		Pos: rl.Vector2{
			X: 5,
			Y: 0,
		},

		Width:  4,
		Height: 3,
		Doors: []Door{
			{
				X:      0,
				Y:      1,
				Facing: core.FacingLeft,
				GoesToRoom: RoomPos{
					RoomId: 0,
					X:      3,
					Y:      1,
				},
			},
			{
				X:      3,
				Y:      1,
				Facing: core.FacingRight,
				GoesToRoom: RoomPos{
					RoomId: 2,
					X:      0,
					Y:      1,
				},
			},
		},
	})

	train.addRoom(Room{
		Pos: rl.Vector2{
			X: 10,
			Y: 0,
		},

		Width:  2,
		Height: 3,
		Doors: []Door{
			{
				X:      0,
				Y:      1,
				Facing: core.FacingLeft,
				GoesToRoom: RoomPos{
					RoomId: 1,
					X:      3,
					Y:      1,
				},
			},
			{
				X:      1,
				Y:      1,
				Facing: core.FacingRight,
				GoesToRoom: RoomPos{
					RoomId: 3,
					X:      0,
					Y:      1,
				},
			},
		},
	})

	train.addRoom(Room{
		Pos: rl.Vector2{
			X: 13,
			Y: 0,
		},

		Width:  2,
		Height: 3,
		Doors: []Door{
			{
				X:      0,
				Y:      1,
				Facing: core.FacingLeft,
				GoesToRoom: RoomPos{
					RoomId: 2,
					X:      3,
					Y:      1,
				},
			},
			{
				X:      1,
				Y:      1,
				Facing: core.FacingRight,
				GoesToRoom: RoomPos{
					RoomId: 3,
					X:      0,
					Y:      1,
				},
			},
		},
	})

	train.addRoom(Room{
		Pos: rl.Vector2{
			X: 16,
			Y: 0,
		},

		Width:  3,
		Height: 3,
		Doors: []Door{
			{
				X:      0,
				Y:      1,
				Facing: core.FacingLeft,
				GoesToRoom: RoomPos{
					RoomId: 3,
					X:      1,
					Y:      1,
				},
			},
		},
	})

	train.addCharacter(NewCharacter("John", 80, RoomPos{RoomId: 0, X: 0, Y: 0}))
	train.addCharacter(NewCharacter("Mary", 60, RoomPos{RoomId: 1, X: 0, Y: 0}))

	return train
}

func (t *Train) addRoom(room Room) {
	room.Id = len(t.Rooms)
	t.Rooms = append(t.Rooms, room)
}

func (t *Train) GetRoom(roomId int) *Room {
	return &t.Rooms[roomId]
}

func (t *Train) FindRoomPath(startRoomId, endRoomId int) ([]int, bool) {
	if startRoomId == endRoomId {
		return []int{startRoomId}, true
	}

	type searchNode struct {
		RoomId int
		Path   []int
	}

	visited := make(map[int]bool, len(t.Rooms))
	queue := []searchNode{{RoomId: startRoomId, Path: []int{startRoomId}}}
	visited[startRoomId] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, door := range t.Rooms[current.RoomId].Doors {
			nextRoomId := door.GoesToRoom.RoomId
			if visited[nextRoomId] {
				continue
			}

			nextPath := append(append([]int{}, current.Path...), nextRoomId)
			if nextRoomId == endRoomId {
				return nextPath, true
			}

			visited[nextRoomId] = true
			queue = append(queue, searchNode{RoomId: nextRoomId, Path: nextPath})
		}
	}

	return nil, false
}

func (t *Train) MoveCharacter(character *Character, target RoomPos) ([]int, error) {
	if character == nil {
		return nil, fmt.Errorf("no character selected")
	}
	if target.RoomId < 0 || target.RoomId >= len(t.Rooms) {
		return nil, fmt.Errorf("invalid room")
	}

	targetRoom := t.GetRoom(target.RoomId)
	if !targetRoom.HasTile(target.X, target.Y) {
		return nil, fmt.Errorf("invalid destination tile")
	}

	if blockingCharacter, ok := t.GetCharacterAtRoomPos(target, character); ok {
		return nil, fmt.Errorf("%s is already occupying %s(%d,%d)", blockingCharacter.Name, string(targetRoom.GetRune()), target.X+1, target.Y+1)
	}

	if character.Pos.RoomId < 0 || character.Pos.RoomId >= len(t.Rooms) {
		return nil, fmt.Errorf("character is not in a valid room")
	}

	path, ok := t.FindRoomPath(character.Pos.RoomId, target.RoomId)
	if !ok {
		return nil, fmt.Errorf("no path found to %s", strings.ToLower(string(targetRoom.GetRune())))
	}

	// Build full animation path with room transitions
	animationPath := make([]RoomPos, 0)
	animationPath = append(animationPath, character.Pos) // Start position

	for i := 0; i < len(path)-1; i++ {
		currentRoomId := path[i]
		nextRoomId := path[i+1]

		currentRoom := t.GetRoom(currentRoomId)
		door := currentRoom.GetDoorTo(nextRoomId)
		if door == nil {
			return nil, fmt.Errorf("no door found between rooms")
		}

		// Move to the door in current room
		animationPath = append(animationPath, RoomPos{RoomId: currentRoomId, X: door.X, Y: door.Y})
		// Enter from door in next room
		animationPath = append(animationPath, RoomPos{RoomId: nextRoomId, X: door.GoesToRoom.X, Y: door.GoesToRoom.Y})
	}

	// Add final target if it's different from last position
	if len(animationPath) == 0 || animationPath[len(animationPath)-1] != target {
		animationPath = append(animationPath, target)
	}

	for i := 1; i < len(animationPath); i++ {
		if blockingCharacter, ok := t.GetCharacterAtRoomPos(animationPath[i], character); ok {
			blockingRoom := t.GetRoom(animationPath[i].RoomId)
			return nil, fmt.Errorf("%s is already occupying %s(%d,%d)", blockingCharacter.Name, string(blockingRoom.GetRune()), animationPath[i].X+1, animationPath[i].Y+1)
		}
	}

	// Set up animation state
	character.IsMoving = true
	character.MovementPath = animationPath
	character.CurrentPathIndex = 0
	character.AnimationProgress = 0.0

	return path, nil
}

func (t *Train) GetCharacterAtRoomPos(target RoomPos, ignore *Character) (*Character, bool) {
	for _, otherCharacter := range t.Characters {
		if otherCharacter == ignore {
			continue
		}

		animatedPos := otherCharacter.GetAnimatedPosition()
		if animatedPos == target {
			return otherCharacter, true
		}
	}

	return nil, false
}

// UpdateCharacterAnimations advances all character animations based on delta time
func (t *Train) UpdateCharacterAnimations(deltaTime float32) {
	for _, character := range t.Characters {
		if !character.IsMoving {
			continue
		}

		// Advance animation progress
		// Calculate distance to next position in tiles
		const distancePerStep float32 = 1.0
		character.AnimationProgress += character.AnimationSpeed * deltaTime / distancePerStep

		if character.AnimationProgress >= 1.0 {
			character.AnimationProgress = 0.0
			character.CurrentPathIndex++

			// Check if animation is complete
			if character.CurrentPathIndex >= len(character.MovementPath) {
				// Set final position and stop animation
				character.Pos = character.MovementPath[len(character.MovementPath)-1]
				character.IsMoving = false
				character.MovementPath = make([]RoomPos, 0)
			}
		}
	}
}

// GetAnimatedPosition returns the current display position of a character (interpolated during movement)
func (character *Character) GetAnimatedPosition() RoomPos {
	if !character.IsMoving || len(character.MovementPath) == 0 {
		return character.Pos
	}

	if character.CurrentPathIndex >= len(character.MovementPath)-1 {
		return character.MovementPath[len(character.MovementPath)-1]
	}

	currentPos := character.MovementPath[character.CurrentPathIndex]
	nextPos := character.MovementPath[character.CurrentPathIndex+1]

	// If crossing between rooms, instantly teleport to the next room
	if currentPos.RoomId != nextPos.RoomId {
		return nextPos
	}

	// Move orthogonally (X first, then Y) to avoid diagonal movement
	moveX := nextPos.X - currentPos.X
	moveY := nextPos.Y - currentPos.Y

	// Determine which axis to move along first
	animX := currentPos.X
	animY := currentPos.Y

	if moveX != 0 {
		// Move along X axis first
		if character.AnimationProgress < 0.5 {
			// First half: move along X
			xProgress := character.AnimationProgress * 2.0 // 0.0 to 1.0 for this half
			animX = currentPos.X + int(float32(moveX)*xProgress)
		} else {
			// Second half: reached target X
			animX = nextPos.X
			// Move along Y in second half
			if moveY != 0 {
				yProgress := (character.AnimationProgress - 0.5) * 2.0 // 0.0 to 1.0 for second half
				animY = currentPos.Y + int(float32(moveY)*yProgress)
			} else {
				animY = nextPos.Y
			}
		}
	} else if moveY != 0 {
		// Only Y moves
		animY = currentPos.Y + int(float32(moveY)*character.AnimationProgress)
	}

	return RoomPos{
		RoomId: currentPos.RoomId,
		X:      animX,
		Y:      animY,
	}
}

func (t *Train) addCharacter(character *Character) {
	character.Id = len(t.Characters)
	t.Characters = append(t.Characters, character)
}
