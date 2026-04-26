package game

import (
	"fmt"
	"math"
)

type WeaponType string

const (
	WeaponCannon  WeaponType = "Cannon"
	WeaponMissile WeaponType = "Missile"
)

type Weapon struct {
	Name              string
	Type              WeaponType
	Damage            int
	CooldownSeconds   float32
	CooldownRemaining float32
	BypassesShields   bool
}

func NewCannon(name string) Weapon {
	return Weapon{
		Name:            name,
		Type:            WeaponCannon,
		Damage:          2,
		CooldownSeconds: 2.0,
	}
}

func NewMissile(name string) Weapon {
	return Weapon{
		Name:            name,
		Type:            WeaponMissile,
		Damage:          3,
		CooldownSeconds: 5.0,
		BypassesShields: true,
	}
}

func (w *Weapon) Ready() bool {
	return w.CooldownRemaining <= 0
}

func (w *Weapon) StartCooldown() {
	w.CooldownRemaining = w.CooldownSeconds
}

func (w *Weapon) AdvanceCooldown(deltaSeconds float32) {
	if w.CooldownRemaining <= 0 || deltaSeconds <= 0 {
		return
	}

	w.CooldownRemaining -= deltaSeconds
	if w.CooldownRemaining < 0 {
		w.CooldownRemaining = 0
	}
}

func (w Weapon) CooldownDisplaySeconds() int {
	return int(math.Ceil(float64(w.CooldownRemaining)))
}

func (w Weapon) ChargeProgress() float32 {
	if w.CooldownSeconds <= 0 {
		return 1
	}

	progress := (w.CooldownSeconds - w.CooldownRemaining) / w.CooldownSeconds
	if progress < 0 {
		return 0
	}
	if progress > 1 {
		return 1
	}

	return progress
}

func (w Weapon) StatusText() string {
	if w.Ready() {
		return fmt.Sprintf("%s ready", w.Name)
	}

	return fmt.Sprintf("%s cd:%ds", w.Name, w.CooldownDisplaySeconds())
}
