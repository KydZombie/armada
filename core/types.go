package core

type Screen interface {
	Update(gm *GameManager)
	Draw(gm *GameManager)
	DrawUI(gm *GameManager)
}
