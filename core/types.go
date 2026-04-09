package core

type Screen interface {
	ResizeScreen(gm *GameManager)

	UpdateScreen(gm *GameManager)
	DrawScreen(gm *GameManager)
	DrawScreenUI(gm *GameManager)
}
