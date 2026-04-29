package core

import rl "github.com/gen2brain/raylib-go/raylib"

type Window[State any] interface {
	IsVisible() bool
	SetVisible(visible bool)

	UpdateWindowSize(gm *GameManager)

	// HandleInput returns true if the input was captured
	HandleInput(gm *GameManager, state *State) bool
	UpdateWindow(gm *GameManager, state *State)
	DrawWindow(gm *GameManager, state *State)
}

type BaseWindow[State any] struct {
	sizeFunc func(gm *GameManager) rl.Rectangle
	bounds   rl.Rectangle
	visible  bool
}

func NewBaseWindow[State any](sizeFunc func(gm *GameManager) rl.Rectangle, gm *GameManager, visible bool) BaseWindow[State] {
	window := BaseWindow[State]{
		sizeFunc: sizeFunc,
		visible:  visible,
	}
	window.UpdateWindowSize(gm)

	return window
}

func (w *BaseWindow[State]) UpdateWindowSize(gm *GameManager) {
	w.bounds = w.sizeFunc(gm)
}

func (w *BaseWindow[State]) IsVisible() bool {
	return w.visible
}

func (w *BaseWindow[State]) SetVisible(visible bool) {
	w.visible = visible
}

func (w *BaseWindow[State]) GetBounds() rl.Rectangle {
	return w.bounds
}

func (w *BaseWindow[State]) SetBounds(bounds rl.Rectangle) {
	w.bounds = bounds
}

// GetTranslatedMousePos returns nil if mouse is not within bounds of window
func (w *BaseWindow[State]) GetTranslatedMousePos() *rl.Vector2 {
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), w.bounds) {
		return new(rl.Vector2Add(rl.GetMousePosition(), rl.Vector2{X: w.bounds.X, Y: w.bounds.Y}))
	}

	return nil
}
