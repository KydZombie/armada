package game

import (
	"math/rand"

	"github.com/KydZombie/armada/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Room struct {
	Id int

	Pos           rl.Vector2
	Width, Height int
	Health        int
	MaxHealth     int
	AttackPower   int

	Doors []Door

	System ShipSystem
}

func (r *Room) ApplyDamage(amount int) int {
	if amount <= 0 {
		return 0
	}

	r.Health -= amount
	if r.Health < 0 {
		r.Health = 0
	}

	return amount
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
	Weapons           []Weapon

	shieldLayers            int
	shieldRechargeRemaining float32
}

func NewTrain(health int) *Train {
	train := &Train{
		Health:     health,
		MaxHealth:  health,
		Rooms:      make([]Room, 0),
		Characters: make([]*Character, 0),
		Weapons:    []Weapon{NewCannon("cannon"), NewMissile("missile")},
	}

	train.addRoom(Room{
		Pos: rl.Vector2{
			X: 0,
			Y: 0,
		},

		Width:       4,
		Height:      3,
		Health:      6,
		MaxHealth:   6,
		AttackPower: 1,
		System:      ShipSystem{Type: SystemPiloting},
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

		Width:       4,
		Height:      3,
		Health:      5,
		MaxHealth:   5,
		AttackPower: 2,
		System:      ShipSystem{Type: SystemWeapons},
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

		Width:       2,
		Height:      3,
		Health:      4,
		MaxHealth:   4,
		AttackPower: 1,
		System:      ShipSystem{Type: SystemEngines},
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

		Width:       2,
		Height:      3,
		Health:      4,
		MaxHealth:   4,
		AttackPower: 1,
		System:      ShipSystem{Type: SystemShields},
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
					RoomId: 4,
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

		Width:       3,
		Height:      3,
		Health:      7,
		MaxHealth:   7,
		AttackPower: 2,
		System:      ShipSystem{Type: SystemMedbay},
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
			{
				X:      2,
				Y:      1,
				Facing: core.FacingRight,
				GoesToRoom: RoomPos{
					RoomId: 5,
					X:      0,
					Y:      1,
				},
			},
		},
	})

	train.addRoom(Room{
		Pos: rl.Vector2{
			X: 19,
			Y: 0,
		},

		Width:       2,
		Height:      3,
		Health:      4,
		MaxHealth:   4,
		AttackPower: 0,
		System:      ShipSystem{Type: SystemLifeSupport},
		Doors: []Door{
			{
				X:      0,
				Y:      1,
				Facing: core.FacingLeft,
				GoesToRoom: RoomPos{
					RoomId: 4,
					X:      2,
					Y:      1,
				},
			},
		},
	})

	train.addCharacter(NewCharacter("John", 80, RoomPos{RoomId: 0, X: 0, Y: 0}))
	train.addCharacter(NewCharacter("Mary", 60, RoomPos{RoomId: 1, X: 0, Y: 0}))
	train.shieldLayers = train.maxShieldLayers()

	return train
}

func (t *Train) addRoom(room Room) {
	room.Id = len(t.Rooms)
	t.Rooms = append(t.Rooms, room)
}

func (t *Train) GetRoom(roomId int) (*Room, bool) {
	if roomId < 0 || roomId >= len(t.Rooms) {
		return nil, false
	}

	return &t.Rooms[roomId], true
}

func (t *Train) addCharacter(character *Character) {
	character.Id = len(t.Characters)
	t.Characters = append(t.Characters, character)
}

func (r *Room) OperationalRatio() float32 {
	if r.MaxHealth <= 0 {
		return 0
	}

	if r.Health <= 0 {
		return 0
	}

	return float32(r.Health) / float32(r.MaxHealth)
}

func (r *Room) IsOperational() bool {
	return r.Health > 0
}

func (t *Train) TotalAttackPower() int {
	total := 0
	if !t.WeaponsOperational() {
		return 0
	}

	for _, weapon := range t.Weapons {
		if weapon.Ready() {
			total += weapon.Damage
		}
	}

	return total
}

func (t *Train) UpdateCombatState(deltaSeconds float32) {
	if deltaSeconds <= 0 {
		return
	}

	maxLayers := t.maxShieldLayers()
	if t.shieldLayers > maxLayers {
		t.shieldLayers = maxLayers
	}

	if t.shieldLayers >= maxLayers {
		t.shieldRechargeRemaining = 0
		return
	}

	if maxLayers <= 0 {
		t.shieldRechargeRemaining = 0
		return
	}

	if t.shieldRechargeRemaining <= 0 {
		t.shieldRechargeRemaining = 4
	}

	t.shieldRechargeRemaining -= deltaSeconds
	for t.shieldRechargeRemaining <= 0 && t.shieldLayers < maxLayers {
		t.shieldLayers++
		if t.shieldLayers >= maxLayers {
			t.shieldRechargeRemaining = 0
			return
		}
		t.shieldRechargeRemaining += 4
	}
}

func (t *Train) WeaponsOperational() bool {
	for _, room := range t.Rooms {
		if room.System.Type == SystemWeapons && room.IsOperational() {
			return true
		}
	}

	return false
}

func (t *Train) AdvanceWeaponCooldowns(deltaSeconds float32) {
	for i := range t.Weapons {
		t.Weapons[i].AdvanceCooldown(deltaSeconds)
	}
}

func (t *Train) GetWeaponByName(name string) (*Weapon, bool) {
	for i := range t.Weapons {
		if t.Weapons[i].Name == name {
			return &t.Weapons[i], true
		}
	}

	return nil, false
}

func (t *Train) ReadyWeapons() int {
	ready := 0
	for _, weapon := range t.Weapons {
		if weapon.Ready() {
			ready++
		}
	}

	return ready
}

func (t *Train) ApplyRoomDamage(roomID int, amount int) int {
	if roomID < 0 || roomID >= len(t.Rooms) {
		return 0
	}

	return t.Rooms[roomID].ApplyDamage(amount)
}

func (t *Train) maxShieldLayers() int {
	layers := 0
	for _, room := range t.Rooms {
		if room.System.Type != SystemShields || !room.IsOperational() {
			continue
		}

		if room.OperationalRatio() >= 0.5 {
			layers++
		}
	}

	return layers
}

func (t *Train) ShieldLayers() int {
	maxLayers := t.maxShieldLayers()
	if t.shieldLayers > maxLayers {
		return maxLayers
	}

	return t.shieldLayers
}

func (t *Train) EvasionChance() int {
	enginesOnline := false
	pilotingOnline := false

	for _, room := range t.Rooms {
		if !room.IsOperational() {
			continue
		}

		switch room.System.Type {
		case SystemEngines:
			enginesOnline = true
		case SystemPiloting:
			pilotingOnline = true
		}
	}

	switch {
	case enginesOnline && pilotingOnline:
		return 20
	case enginesOnline || pilotingOnline:
		return 10
	default:
		return 0
	}
}

func (t *Train) MedbayHealingPerTick() int {
	for _, room := range t.Rooms {
		if room.System.Type == SystemMedbay && room.IsOperational() {
			if room.OperationalRatio() >= 0.5 {
				return 2
			}
			return 1
		}
	}

	return 0
}

func (t *Train) LifeSupportOperational() bool {
	for _, room := range t.Rooms {
		if room.System.Type == SystemLifeSupport && room.IsOperational() {
			return true
		}
	}

	return false
}

func (t *Train) LifeSupportDamagePerTick() int {
	for _, room := range t.Rooms {
		if room.System.Type != SystemLifeSupport {
			continue
		}

		if !room.IsOperational() {
			return 2
		}

		if room.OperationalRatio() < 0.5 {
			return 1
		}

		return 0
	}

	return 0
}

func (t *Train) ApplyHullDamage(amount int) int {
	if amount <= 0 {
		return 0
	}

	t.Health -= amount
	if t.Health < 0 {
		t.Health = 0
	}

	return amount
}

func (t *Train) ResolveIncomingAttack(amount int) (hullDamage int, evaded bool, shielded bool) {
	if amount <= 0 {
		return 0, false, true
	}

	if chance := t.EvasionChance(); chance > 0 && rand.Intn(100) < chance {
		return 0, true, false
	}

	if t.ShieldLayers() > 0 {
		t.shieldLayers--
		if t.shieldRechargeRemaining <= 0 {
			t.shieldRechargeRemaining = 4
		}
		return 0, false, true
	}

	return t.ApplyHullDamage(amount), false, false
}
