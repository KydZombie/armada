package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func drawStatBar(bounds rl.Rectangle, label string, value int, maxValue int, fillColor rl.Color, textSize int32) {
	if maxValue <= 0 {
		maxValue = 1
	}

	if value < 0 {
		value = 0
	}

	if value > maxValue {
		value = maxValue
	}

	rl.DrawRectangleRec(bounds, rl.DarkGray)

	fillWidth := bounds.Width * float32(value) / float32(maxValue)
	if fillWidth > 0 {
		rl.DrawRectangleRec(rl.Rectangle{
			X:      bounds.X,
			Y:      bounds.Y,
			Width:  fillWidth,
			Height: bounds.Height,
		}, fillColor)
	}

	rl.DrawRectangleLinesEx(bounds, 2, rl.White)

	labelText := fmt.Sprintf("%s: %d/%d", label, value, maxValue)
	textWidth := rl.MeasureText(labelText, textSize)
	textX := int32(bounds.X + (bounds.Width-float32(textWidth))/2)
	textY := int32(bounds.Y + (bounds.Height-float32(textSize))/2)
	rl.DrawText(labelText, textX, textY, textSize, rl.White)
}
