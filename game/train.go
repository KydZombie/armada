package game

type Train struct {
	Health, MaxHealth int
}

func NewTrain() *Train {
	return &Train{
		Health:    100,
		MaxHealth: 100,
	}
}
