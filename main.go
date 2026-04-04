package main

import (
	"time"

	"github.com/gen2brain/raylib-go/raylib"
)

func initWindow() {
	//rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(800, 450, "Armada")
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())))
}

func main() {
	initWindow()
	defer rl.CloseWindow()

	game := Game{}
	game.init()
	defer game.close()

	lastTime := time.Now().UnixMilli()

	for !rl.WindowShouldClose() {
		currentTime := time.Now().UnixMilli()
		delta := float32(currentTime-lastTime) / 1000.0
		lastTime = currentTime

		game.update(delta)

		rl.BeginDrawing()
		game.draw()
		game.gui()
		rl.DrawFPS(0, 0)

		rl.EndDrawing()
	}

}
