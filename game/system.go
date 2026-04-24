package game

type SystemType string

const (
	SystemWeapons     SystemType = "Weapons"
	SystemEngines     SystemType = "Engines"
	SystemPiloting    SystemType = "Piloting"
	SystemShields     SystemType = "Shields"
	SystemMedbay      SystemType = "Medbay"
	SystemLifeSupport SystemType = "Life Support"
)

type ShipSystem struct {
	Type SystemType
}

func (s ShipSystem) Name() string {
	return string(s.Type)
}

func (s ShipSystem) ShortName() string {
	switch s.Type {
	case SystemWeapons:
		return "WPN"
	case SystemEngines:
		return "ENG"
	case SystemPiloting:
		return "PIL"
	case SystemShields:
		return "SHD"
	case SystemMedbay:
		return "MED"
	case SystemLifeSupport:
		return "LIF"
	default:
		return "SYS"
	}
}
