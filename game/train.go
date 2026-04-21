package game

import (
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
			X: 12,
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
			X: 15,
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

func (t *Train) addCharacter(character *Character) {
	character.Id = len(t.Characters)
	t.Characters = append(t.Characters, character)
}
