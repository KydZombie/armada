package core

type Screen interface {
	UpdateScreen(gm *GameManager)
	DrawScreen(gm *GameManager)
	DrawScreenUI(gm *GameManager)
}
