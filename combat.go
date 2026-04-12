package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Direction int

const (
	// DirectionDown is the default facing when the player is standing still
	// and no more recent movement input has changed their orientation.
	DirectionDown Direction = iota
	// DirectionUp means the player is facing toward the top of the screen.
	DirectionUp
	// DirectionLeft means the player is facing toward the left side of the screen.
	DirectionLeft
	// DirectionRight means the player is facing toward the right side of the screen.
	DirectionRight
)

const (
	// The constants define the fixed size of the test combat window.
	screenWidth  int32 = 1000
	screenHeight int32 = 700

	// windowTitle is used when opening the standalone combat demo window.
	windowTitle = "Armada - Enemy Test"
)

// Player stores player position, size, speed, and attack timing.
type Player struct {
	X               float32
	Y               float32
	Width           float32
	Height          float32
	Speed           float32
	AttackReach     float32
	AttackThickness float32
	AttackCooldown  float32
	LastAttackTime  float64
	Facing          Direction
}

// Enemy stores enemy position, size, health, and alive state.
type Enemy struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
	Health int
	Alive  bool
}

// GetRect returns the player's rectangle for drawing/collision.
func (p Player) GetRect() rl.Rectangle {
	return rl.Rectangle{
		X:      p.X,
		Y:      p.Y,
		Width:  p.Width,
		Height: p.Height,
	}
}

// Center returns the midpoint of the player rectangle.
// This is useful for positioning directional attack hitboxes so they line up
// with the middle of the player sprite instead of the player's top-left corner.
func (p Player) Center() rl.Vector2 {
	return rl.Vector2{
		X: p.X + p.Width/2,
		Y: p.Y + p.Height/2,
	}
}

// GetRect returns the enemy's rectangle for drawing/collision.
func (e Enemy) GetRect() rl.Rectangle {
	return rl.Rectangle{
		X:      e.X,
		Y:      e.Y,
		Width:  e.Width,
		Height: e.Height,
	}
}

// GetAttackHitbox returns a short rectangular strike zone in front of the player.
// The hitbox changes shape and position based on the player's current facing:
// vertical attacks are tall and narrow, while horizontal attacks are wide and short.
func (p Player) GetAttackHitbox() rl.Rectangle {
	center := p.Center()

	switch p.Facing {
	case DirectionUp:
		return rl.Rectangle{
			X:      center.X - p.AttackThickness/2,
			Y:      p.Y - p.AttackReach,
			Width:  p.AttackThickness,
			Height: p.AttackReach,
		}
	case DirectionLeft:
		return rl.Rectangle{
			X:      p.X - p.AttackReach,
			Y:      center.Y - p.AttackThickness/2,
			Width:  p.AttackReach,
			Height: p.AttackThickness,
		}
	case DirectionRight:
		return rl.Rectangle{
			X:      p.X + p.Width,
			Y:      center.Y - p.AttackThickness/2,
			Width:  p.AttackReach,
			Height: p.AttackThickness,
		}
	default:
		return rl.Rectangle{
			X:      center.X - p.AttackThickness/2,
			Y:      p.Y + p.Height,
			Width:  p.AttackThickness,
			Height: p.AttackReach,
		}
	}
}

// String returns a readable version of the direction for debug text on screen.
func (d Direction) String() string {
	switch d {
	case DirectionUp:
		return "Up"
	case DirectionLeft:
		return "Left"
	case DirectionRight:
		return "Right"
	default:
		return "Down"
	}
}

// CanAttack checks whether enough time has passed since the last attack.
func (p Player) CanAttack(currentTime float64) bool {
	return currentTime-p.LastAttackTime >= float64(p.AttackCooldown)
}

// getMovementInput reads keyboard input and converts it into per-frame movement.
// The returned values are already scaled by player speed.
//
// When moving diagonally, the vector is normalized so diagonal movement is not
// faster than moving in a straight line.
func getMovementInput(speed float32) (float32, float32) {
	moveX := float32(0)
	moveY := float32(0)

	if rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD) {
		moveX++
	}
	if rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA) {
		moveX--
	}
	if rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(rl.KeyW) {
		moveY--
	}
	if rl.IsKeyDown(rl.KeyDown) || rl.IsKeyDown(rl.KeyS) {
		moveY++
	}

	if moveX != 0 && moveY != 0 {
		scale := float32(1 / math.Sqrt2)
		moveX *= scale
		moveY *= scale
	}

	return moveX * speed, moveY * speed
}

// updateFacing stores the most recent movement direction on the player.
// The combat system uses this facing to decide where the next attack hitbox
// should appear. Horizontal movement is prioritized when both axes are active.
func updateFacing(player *Player, moveX, moveY float32) {
	if moveX > 0 {
		player.Facing = DirectionRight
		return
	}
	if moveX < 0 {
		player.Facing = DirectionLeft
		return
	}
	if moveY < 0 {
		player.Facing = DirectionUp
		return
	}
	if moveY > 0 {
		player.Facing = DirectionDown
	}
}

// clampPlayerToScreen keeps the player's rectangle fully inside the window.
// This prevents movement and attack origin points from drifting off-screen.
func clampPlayerToScreen(player *Player, width, height int32) {
	if player.X < 0 {
		player.X = 0
	}
	if player.Y < 0 {
		player.Y = 0
	}
	if player.X+player.Width > float32(width) {
		player.X = float32(width) - player.Width
	}
	if player.Y+player.Height > float32(height) {
		player.Y = float32(height) - player.Height
	}
}

// handleAttack processes a single attack press for the current frame.
// It returns:
// - whether an attack actually happened this frame
// - the message that should be shown to the player
//
// The attack only triggers when SPACE is pressed and the cooldown has elapsed.
// If the enemy overlaps the current directional hitbox, the enemy takes damage.
func handleAttack(player *Player, enemy *Enemy, currentTime float64) (bool, string) {
	if !rl.IsKeyPressed(rl.KeySpace) || !player.CanAttack(currentTime) {
		return false, ""
	}

	player.LastAttackTime = currentTime
	attackRect := player.GetAttackHitbox()

	if enemy.Alive && rl.CheckCollisionRecs(attackRect, enemy.GetRect()) {
		enemy.Health--
		if enemy.Health <= 0 {
			enemy.Alive = false
			return true, "Enemy defeated!"
		}

		return true, "Hit enemy!"
	}

	return true, "Attack missed!"
}

// newPlayer creates the default player state used by the demo.
// Keeping setup in one place makes balancing and later reuse easier.
func newPlayer() Player {
	return Player{
		X:               120,
		Y:               250,
		Width:           50,
		Height:          50,
		Speed:           4,
		AttackReach:     80,
		AttackThickness: 42,
		AttackCooldown:  0.4, // seconds
		LastAttackTime:  -1,
		Facing:          DirectionRight,
	}
}

// newEnemy creates the single target dummy used in this prototype.
func newEnemy() Enemy {
	return Enemy{
		X:      500,
		Y:      260,
		Width:  60,
		Height: 60,
		Health: 3,
		Alive:  true,
	}
}

// drawScene renders the entire combat test frame.
// The update logic is kept outside this function so drawing stays predictable
// and focused only on presentation.
func drawScene(player Player, enemy Enemy, lastMessage string, showAttackBox bool) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	// Background test panels
	rl.DrawRectangle(40, 80, 420, 540, rl.SkyBlue)
	rl.DrawRectangle(520, 80, 420, 540, rl.Pink)

	// Player
	rl.DrawRectangleRec(player.GetRect(), rl.Blue)
	rl.DrawText("PLAYER", int32(player.X), int32(player.Y)-22, 18, rl.DarkBlue)
	rl.DrawText(
		fmt.Sprintf("Facing: %s", player.Facing.String()),
		int32(player.X)-8,
		int32(player.Y)+int32(player.Height)+8,
		18,
		rl.DarkBlue,
	)

	// Enemy
	if enemy.Alive {
		rl.DrawRectangleRec(enemy.GetRect(), rl.Red)
		rl.DrawText(
			fmt.Sprintf("Enemy HP: %d", enemy.Health),
			int32(enemy.X)-10,
			int32(enemy.Y)-25,
			20,
			rl.Maroon,
		)
	} else {
		rl.DrawText("Enemy defeated", int32(enemy.X)-10, int32(enemy.Y)+20, 20, rl.Gray)
	}

	// Show attack box only on the frame attack happens
	if showAttackBox {
		rl.DrawRectangleLinesEx(player.GetAttackHitbox(), 3, rl.Gold)
		rl.DrawText("ATTACK", int32(player.X)-5, int32(player.Y)-45, 18, rl.Gold)
	}

	// Instructions
	rl.DrawText("Move: WASD or Arrow Keys", 30, 20, 24, rl.Black)
	rl.DrawText("Attack: SPACE", 30, 50, 24, rl.Black)
	rl.DrawText(lastMessage, 30, 650, 24, rl.DarkGray)

	rl.EndDrawing()
}

// RunCombatDemo starts the standalone combat prototype loop.
// Keep this separate from the package entrypoint so the file can live alongside the main game.
func RunCombatDemo() {
	rl.InitWindow(screenWidth, screenHeight, windowTitle)
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	player := newPlayer()
	enemy := newEnemy()

	lastMessage := "Press SPACE to attack"
	showAttackBox := false

	for !rl.WindowShouldClose() {
		// Raylib provides elapsed time in seconds, which we use for attack cooldowns.
		currentTime := rl.GetTime()

		// This visual is only shown on the frame where the attack is triggered.
		showAttackBox = false

		// Collect movement input, apply it, then update the player's last facing
		// direction so attacks follow the direction of travel.
		moveX, moveY := getMovementInput(player.Speed)
		player.X += moveX
		player.Y += moveY
		updateFacing(&player, moveX, moveY)

		// Make sure the player stays inside the test arena.
		clampPlayerToScreen(&player, screenWidth, screenHeight)

		// If an attack happened this frame, update the on-screen feedback.
		if attacked, message := handleAttack(&player, &enemy, currentTime); attacked {
			showAttackBox = true
			lastMessage = message
		}

		// Draw the current world state after all updates are complete.
		drawScene(player, enemy, lastMessage, showAttackBox)
	}
}
