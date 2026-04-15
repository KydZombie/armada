package game

import "github.com/KydZombie/armada/core"

type Room struct {
	Id int

	X, Y          int
	Width, Height int

	Doors []Door

	System System
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
}

func NewTrain(health int) *Train {
	train := &Train{
		Health:    health,
		MaxHealth: health,
		Rooms:     make([]Room, 0),
	}

	train.addRoom(Room{
		X: 0,
		Y: 0,

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
		X: 5,
		Y: 0,

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
		X: 10,
		Y: 0,

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
		X: 12,
		Y: 0,

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
		X: 15,
		Y: 0,

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

	return train
}

func (t *Train) addRoom(room Room) {
	room.Id = len(room.Doors)
	t.Rooms = append(t.Rooms, room)
}
