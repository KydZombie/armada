package util

import rl "github.com/gen2brain/raylib-go/raylib"

func KeyToAlphanumeric(key int32) (rune, bool) {
	// Only handle printable keys (A-Z, 0-9, etc.)
	switch {
	case key >= rl.KeyA && key <= rl.KeyZ:
		// Convert KeyA..KeyZ to 'a'..'z'
		return 'a' + (key - rl.KeyA), true
	case key >= rl.KeyZero && key <= rl.KeyNine:
		return '0' + (key - rl.KeyZero), true
	case key == rl.KeyMinus:
		if rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift) {
			return '_', true
		}
		return '-', true
	case key == rl.KeySpace:
		return ' ', true
	default:
		return 0, false // Not a handled key
	}
}
