package game

import "github.com/KydZombie/armada/core"

type Character struct {
	Id int

	Name              string
	Health, MaxHealth int

	Pos    RoomPos
	Facing core.Facing
}

func NewCharacter(name string, health int, pos RoomPos) *Character {
	return &Character{
		Name:      name,
		Health:    health,
		MaxHealth: health,
		Pos:       pos,
		Facing:    core.FacingUp,
	}
}
