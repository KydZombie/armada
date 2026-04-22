package game

import "math/rand"

// Enemy describes the small set of behaviors the rest of the game
// can expect from an enemy. Right now we only need to know whether
// the enemy is still alive, how it takes damage, and how much damage
// it deals when it attacks.
type Enemy interface {
	Name() string
	Health() int
	MaxHealth() int
	Alive() bool
	TakeDamage(amount int)
	Attack() int
	AdvanceTurn() (bool, int)
}

// EnemyPart is groundwork for future targeted combat, such as aiming
// attacks at the head or body. The current prototype does not use
// parts in combat yet, but storing them now makes later expansion easier.
type EnemyPart struct {
	Name string

	Health int

	MaxHealth int
}

// BasicEnemy is a simple starter enemy implementation for the prototype.
// The fields are intentionally small and direct so the behavior is easy
// to understand and expand later.
type BasicEnemy struct {
	// name is the display name of the enemy.
	name string

	// health is the enemy's current health.
	health int

	// maxHealth is the highest health value the enemy can have.
	maxHealth int

	// attackPower is the amount of damage this enemy deals when it attacks.
	attackPower int

	// attackCooldown stores how many turns this enemy waits between attacks.
	attackCooldown int

	// attackTimer stores the current turn countdown until the next attack.
	attackTimer int

	// Parts is groundwork for future targeted attacks. For now, combat still
	// uses the main enemy health values, so these parts are informational only.
	Parts []*EnemyPart
}

// NewBasicEnemy creates a basic enemy with the provided stats.
// health and maxHealth start with the same value because a new enemy
// should begin at full health.
func NewBasicEnemy(name string, health int, attackPower int) *BasicEnemy {
	const defaultAttackCooldown = 2

	return &BasicEnemy{
		name:           name,
		health:         health,
		maxHealth:      health,
		attackPower:    attackPower,
		attackCooldown: defaultAttackCooldown,
		attackTimer:    defaultAttackCooldown,
		Parts: []*EnemyPart{
			{
				Name:      "Core",
				Health:    health,
				MaxHealth: health,
			},
			{
				Name:      "Armor",
				Health:    health,
				MaxHealth: health,
			},
			{
				Name:      "Tail",
				Health:    health,
				MaxHealth: health,
			},
		},
	}
}

// Name returns the enemy's display name.
func (e *BasicEnemy) Name() string {
	return e.name
}

// Health returns the enemy's current health.
func (e *BasicEnemy) Health() int {
	return e.health
}

// MaxHealth returns the enemy's highest possible health value.
func (e *BasicEnemy) MaxHealth() int {
	return e.maxHealth
}

// Alive returns true while the enemy still has health remaining.
func (e *BasicEnemy) Alive() bool {
	return e.health > 0
}

// TakeDamage lowers the enemy's health by the given amount.
// Health is clamped at 0 so it never becomes negative.
func (e *BasicEnemy) TakeDamage(amount int) {
	// Ignore zero or negative damage so this method only handles
	// real damage values.
	if amount <= 0 {
		return
	}

	e.health -= amount

	if e.health < 0 {
		e.health = 0
	}
}

// Attack returns the amount of damage this enemy should deal.
func (e *BasicEnemy) Attack() int {
	return e.attackPower
}

func (e *BasicEnemy) AdvanceTurn() (bool, int) {
	if !e.Alive() || e.attackCooldown <= 0 {
		return false, 0
	}

	if e.attackTimer > 0 {
		e.attackTimer--
	}

	if e.attackTimer > 0 {
		return false, 0
	}

	damage := e.Attack()
	e.attackTimer = randomEnemyCooldown()
	return true, damage
}

func randomEnemyCooldown() int {
	const minEnemyCooldown = 1
	const maxEnemyCooldown = 3

	return rand.Intn(maxEnemyCooldown-minEnemyCooldown+1) + minEnemyCooldown
}
