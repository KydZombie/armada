package game

import (
	"fmt"
	"strings"
)

// Weapon tracks a train weapon's turn-based cooldown.
type Weapon struct {
	Name string

	Damage int

	CooldownTurns     int
	CooldownRemaining int
}

func NewWeapon(name string, damage int, cooldownTurns int) Weapon {
	if cooldownTurns < 0 {
		cooldownTurns = 0
	}

	return Weapon{
		Name:          name,
		Damage:        damage,
		CooldownTurns: cooldownTurns,
	}
}

func (w *Weapon) Ready() bool {
	return w.CooldownRemaining <= 0
}

func (w *Weapon) TriggerCooldown() {
	w.CooldownRemaining = w.CooldownTurns
}

func (w *Weapon) AdvanceTurn() {
	if w.CooldownRemaining > 0 {
		w.CooldownRemaining--
	}
}

// AdvancePlayerTurn advances shared combat time after a successful command.
func (g *Game) AdvancePlayerTurn(cmd string) []string {
	messages := make([]string, 0)

	if cmd == "" {
		return messages
	}

	g.PlayerWeapon.AdvanceTurn()

	for _, enemy := range g.RoomEnemies {
		if enemy == nil || !enemy.Alive() {
			continue
		}

		didAttack, damage := enemy.AdvanceTurn()
		if !didAttack || damage <= 0 {
			continue
		}

		g.Train.Health -= damage
		if g.Train.Health < 0 {
			g.Train.Health = 0
		}

		messages = append(messages, fmt.Sprintf(
			"%s attacks the train for %d damage. Train HP: %d/%d.",
			enemy.Name(),
			damage,
			g.Train.Health,
			g.Train.MaxHealth,
		))
	}

	g.syncSelectedRoom()

	return messages
}

func (g *Game) FireSelectedWeapon(targetPart string) (string, bool) {
	g.syncSelectedRoom()

	if g.Enemy == nil {
		return "There is no enemy to attack.", false
	}

	if !g.Enemy.Alive() {
		return fmt.Sprintf("%s is already defeated.", g.Enemy.Name()), false
	}

	if !g.PlayerWeapon.Ready() {
		return fmt.Sprintf(
			"%s is on cooldown for %d more turn(s).",
			g.PlayerWeapon.Name,
			g.PlayerWeapon.CooldownRemaining,
		), false
	}

	damage := g.PlayerWeapon.Damage
	if damage <= 0 {
		damage = 1
	}

	if targetPart != "" {
		enemy, ok := g.Enemy.(*BasicEnemy)
		if !ok {
			return "Enemy parts are not available.", false
		}

		for _, part := range enemy.Parts {
			if part == nil {
				continue
			}

			if !strings.EqualFold(part.Name, targetPart) {
				continue
			}

			part.Health -= damage
			if part.Health < 0 {
				part.Health = 0
			}

			if strings.EqualFold(part.Name, "Core") {
				g.Enemy.TakeDamage(damage)
			}

			g.PlayerWeapon.TriggerCooldown()

			return fmt.Sprintf(
				"%s's %s takes %d dmg. %s HP: %d/%d",
				g.Enemy.Name(),
				part.Name,
				damage,
				part.Name,
				part.Health,
				part.MaxHealth,
			), true
		}

		return "That part does not exist.", false
	}

	g.Enemy.TakeDamage(damage)
	g.PlayerWeapon.TriggerCooldown()

	if g.Enemy.Alive() {
		return fmt.Sprintf(
			"You attack %s for %d damage. %s has %d/%d health left.",
			g.Enemy.Name(),
			damage,
			g.Enemy.Name(),
			g.Enemy.Health(),
			g.Enemy.MaxHealth(),
		), true
	}

	return fmt.Sprintf(
		"You attack %s for %d damage. %s is defeated.",
		g.Enemy.Name(),
		damage,
		g.Enemy.Name(),
	), true
}
