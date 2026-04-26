package game

import "github.com/KydZombie/armada/core"

type Character struct {
	Id int

	Name              string
	Health, MaxHealth int

	Pos    RoomPos
	Facing core.Facing

	// Movement animation state
	IsMoving          bool
	MovementPath      []RoomPos // Path of positions to animate through
	CurrentPathIndex  int       // Which position in the path we're animating to
	AnimationProgress float32   // 0.0 to 1.0, interpolation between current and next position
	AnimationSpeed    float32   // tiles per second (how fast to move)
}

func NewCharacter(name string, health int, pos RoomPos) *Character {
	return &Character{
		Name:              name,
		Health:            health,
		MaxHealth:         health,
		Pos:               pos,
		Facing:            core.FacingUp,
		IsMoving:          false,
		MovementPath:      make([]RoomPos, 0),
		CurrentPathIndex:  0,
		AnimationProgress: 0.0,
		AnimationSpeed:    0.4, // 0.4 tiles per second
	}
}
