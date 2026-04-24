package game

import (
	"fmt"
	"strconv"

	"github.com/KydZombie/armada/core"
)

func initializeCommands() *core.CommandDB[Game] {
	db := core.NewCommandDB[Game]()

	registerGenericCommands(db)
	registerInfoCommands(db)
	registerUnitCommands(db)
	registerCombatCommands(db)
	registerDebugCommands(db)

	return db
}

func parseRoomIndex(arg string, roomCount int) (int, string, bool) {
	roomRunes := []rune(arg)
	if len(roomRunes) != 1 {
		return 0, "Invalid room", false
	}

	roomIdx := int(roomRunes[0] - 'a')
	if roomIdx < 0 || roomIdx >= roomCount {
		return 0, "Invalid room", false
	}

	return roomIdx, "", true
}

func registerGenericCommands(db *core.CommandDB[Game]) {
	// TODO: Need to be able to do more than just this
	db.RegisterCommand(core.Command[Game]{
		Name: "help",
		OnRun: func(args []string, game *Game) (string, bool) {
			for _, command := range db.CmdMap {
				if command.Description != nil {
					game.Terminal.OutputText(fmt.Sprint(command.Name, " ^"), StandardOutputMessage)
					for _, line := range command.Description {
						game.Terminal.OutputText(fmt.Sprint("  ", line), StandardOutputMessage)
					}
				}
			}

			return "", true
		},
		Description: []string{"This help command!"},
	})
}

func registerInfoCommands(db *core.CommandDB[Game]) {

}

func registerUnitCommands(db *core.CommandDB[Game]) {
	db.RegisterCommand(core.Command[Game]{
		Name: "selc",
		OnRun: func(args []string, game *Game) (string, bool) {
			if len(args) != 1 {
				return "Invalid number of arguments", false
			}

			characterIdx, err := strconv.Atoi(args[0])
			characterIdx -= 1
			if err != nil || characterIdx < 0 || characterIdx >= len(game.Train.Characters) {
				return fmt.Sprintf("Invalid character index (%s).", args[0]), false
			}
			game.SelectedCharacterIndex = characterIdx
			character := game.Train.Characters[characterIdx]

			return fmt.Sprint("Selected ", character.Name, "."), true
		},
		Description: []string{"Select a character"},
	})

	db.RegisterCommand(core.Command[Game]{
		Name: "move",
		OnRun: func(args []string, game *Game) (string, bool) {
			if len(args) != 3 {
				return "Invalid number of arguments", false
			}

			if game.SelectedCharacterIndex < 0 {
				return "No character selected", false
			}

			roomIdx, errText, ok := parseRoomIndex(args[0], len(game.Train.Rooms))
			if !ok {
				return errText, false
			}

			room := game.Train.Rooms[roomIdx]

			x, err := strconv.Atoi(args[1])
			if err != nil || x < 0 || x >= room.Width {
				return "Invalid x", false
			}
			y, err := strconv.Atoi(args[2])
			if err != nil || y < 0 || y >= room.Height {
				return "Invalid y", false
			}

			character := game.Train.Characters[game.SelectedCharacterIndex]
			character.Pos.RoomId = roomIdx
			character.Pos.X = x
			character.Pos.Y = y

			game.SelectedCharacterIndex = -1

			return fmt.Sprint("Moved ", character.Name), true
		},
		Description: []string{"Move a character", "Takes arguments [room], [x], and [y]"},
	})
}

func registerCombatCommands(db *core.CommandDB[Game]) {
	attackEnemyCommand := core.Command[Game]{
		Name: "attack_enemy",
		OnRun: func(args []string, game *Game) (string, bool) {
			if game.Enemy == nil {
				return "There is no enemy to attack.", false
			}

			if !game.Enemy.Alive() {
				return fmt.Sprintf("%s is already defeated.", game.Enemy.Name()), false
			}

			damage := game.Train.TotalAttackPower()
			game.Enemy.TakeDamage(damage)

			if game.Enemy.Alive() {
				return fmt.Sprintf(
					"You attack the %s for %d damage. %s has %d/%d health left.",
					game.Enemy.Name(),
					damage,
					game.Enemy.Name(),
					game.Enemy.Health(),
					game.Enemy.MaxHealth(),
				), true
			}

			return fmt.Sprintf(
				"You attack the %s for %d damage. %s is defeated.",
				game.Enemy.Name(),
				damage,
				game.Enemy.Name(),
			), true
		},
		Description: []string{"Attack the current enemy."},
	}
	db.RegisterCommand(attackEnemyCommand)
	db.RegisterCommand(core.Command[Game]{
		Name:        "attack",
		OnRun:       attackEnemyCommand.OnRun,
		Description: []string{"Alias for attack_enemy."},
	})

	enemyStatusCommand := core.Command[Game]{
		Name: "enemy_status",
		OnRun: func(args []string, game *Game) (string, bool) {
			if game.Enemy == nil {
				return "There is no enemy to inspect.", false
			}

			statusText := "alive"
			if !game.Enemy.Alive() {
				statusText = "dead"
			}

			return fmt.Sprintf(
				"Enemy: %s. HP: %d/%d. Status: %s.",
				game.Enemy.Name(),
				game.Enemy.Health(),
				game.Enemy.MaxHealth(),
				statusText,
			), true
		},
		Description: []string{"Get the status of the enemy."},
	}
	db.RegisterCommand(enemyStatusCommand)
	db.RegisterCommand(core.Command[Game]{
		Name:        "status",
		OnRun:       enemyStatusCommand.OnRun,
		Description: []string{"Alias for enemy_status."},
	})
	db.RegisterCommand(core.Command[Game]{
		Name: "train_status",
		OnRun: func(args []string, game *Game) (string, bool) {
			lifeSupportText := "online"
			if !game.Train.LifeSupportOperational() {
				lifeSupportText = fmt.Sprintf("offline (%d dmg/tick)", game.Train.LifeSupportDamagePerTick())
			}

			return fmt.Sprintf(
				"Train hull: %d/%d. Total damage: %d. Shields: %d. Evasion: %d%%. Medbay: %d HP/tick. Life Support: %s.",
				game.Train.Health,
				game.Train.MaxHealth,
				game.Train.TotalAttackPower(),
				game.Train.ShieldLayers(),
				game.Train.EvasionChance(),
				game.Train.MedbayHealingPerTick(),
				lifeSupportText,
			), true
		},
		Description: []string{"Get the status of the train hull and total damage."},
	})

	db.RegisterCommand(core.Command[Game]{
		Name: "damage_train",
		OnRun: func(args []string, game *Game) (string, bool) {
			if len(args) != 1 {
				return "Invalid number of arguments", false
			}

			damage, err := strconv.Atoi(args[0])
			if err != nil || damage <= 0 {
				return "Invalid damage amount", false
			}

			game.Train.Health -= damage
			if game.Train.Health < 0 {
				game.Train.Health = 0
			}

			return fmt.Sprintf(
				"Train takes %d damage. Hull is now %d/%d.",
				damage,
				game.Train.Health,
				game.Train.MaxHealth,
			), true
		},
		Description: []string{"Damage the train hull.", "Takes arguments [amount]"},
	})

	db.RegisterCommand(core.Command[Game]{
		Name: "damage_room",
		OnRun: func(args []string, game *Game) (string, bool) {
			if len(args) != 2 {
				return "Invalid number of arguments", false
			}

			roomIdx, errText, ok := parseRoomIndex(args[0], len(game.Train.Rooms))
			if !ok {
				return errText, false
			}

			damage, err := strconv.Atoi(args[1])
			if err != nil || damage <= 0 {
				return "Invalid damage amount", false
			}

			room := &game.Train.Rooms[roomIdx]
			room.Health -= damage
			if room.Health < 0 {
				room.Health = 0
			}

			return fmt.Sprintf(
				"Room %c takes %d damage. Room health is now %d/%d.",
				room.GetRune(),
				damage,
				room.Health,
				room.MaxHealth,
			), true
		},
		Description: []string{"Damage a train room.", "Takes arguments [room] and [amount]"},
	})
}

func registerDebugCommands(db *core.CommandDB[Game]) {
}
