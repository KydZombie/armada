package game

import (
	"math/rand"

	"fmt"
	"strings"

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
			X: 13,
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
			X: 16,
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

	targetRoom, _ := t.GetRoom(target.RoomId)
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

		currentRoom, _ := t.GetRoom(currentRoomId)
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
			blockingRoom, _ := t.GetRoom(animationPath[i].RoomId)
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
