package core

import rl "github.com/gen2brain/raylib-go/raylib"

type Window[State any] interface {
	IsVisible() bool
	SetVisible(visible bool)

	// HandleInput returns true if the input was captured
	HandleInput(gm *GameManager, state *State) bool
	UpdateWindow(gm *GameManager, state *State)
	DrawWindow(gm *GameManager, state *State)
	DrawWindowUI(gm *GameManager, state *State)
}

type BaseWindow[State any] struct {
	bounds  rl.Rectangle
	visible bool
}

func NewBaseWindow[State any](bounds rl.Rectangle, visible bool) BaseWindow[State] {
	return BaseWindow[State]{
		bounds:  bounds,
		visible: visible,
	}
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
