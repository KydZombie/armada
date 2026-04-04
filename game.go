package main

import (
	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	moving bool
	pos    rl.Vector2
}

func (g *Game) init() {
	g.moving = false
	g.pos = rl.Vector2{X: 200, Y: 50}
}

func (g *Game) close() {}

func (g *Game) update(delta float32) {
	if g.moving {
		g.pos.X += 100.0 * delta
	}
}

func (g *Game) draw() {
	rl.ClearBackground(rl.DarkGray)
	rl.DrawRectangleV(g.pos, rl.Vector2{X: 50, Y: 50}, rl.Red)
}

func (g *Game) gui() {
	var buttonText string
	if g.moving {
		buttonText = "Stop moving"
	} else {
		buttonText = "Start moving"
	}

	if rg.Button(rl.Rectangle{
		X:      0,
		Y:      50,
		Width:  100,
		Height: 60,
	}, buttonText) {
		g.moving = !g.moving
	}
}
