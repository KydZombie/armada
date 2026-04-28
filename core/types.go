package core

type Facing byte

const (
	FacingLeft Facing = iota
	FacingRight
	FacingUp
	FacingDown
)

type Screen interface {
	ResizeScreen(gm *GameManager)

	UpdateScreen(gm *GameManager)
	DrawScreen(gm *GameManager)
}
