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

			roomRunes := []rune(args[0])
			if len(roomRunes) != 1 {
				return "Invalid room", false
			}

			roomRune := roomRunes[0]
			roomIdx := int(roomRune - 'a')

			if roomIdx < 0 || roomIdx >= len(game.Train.Rooms) {
				return "Invalid room", false
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
	// attack is a very small prototype combat command.
	// It always deals 1 damage to the current enemy so we can test
	// terminal commands and battle window updates together.
	db.RegisterCommand(core.Command[Game]{
		Name: "attack",
		OnRun: func(args []string, game *Game) (string, bool) {
			if game.Enemy == nil {
				return "There is no enemy to attack.", false
			}

			if !game.Enemy.Alive() {
				return fmt.Sprintf("%s is already defeated.", game.Enemy.Name()), false
			}

			damage := 1
			game.Enemy.TakeDamage(damage)

			if game.Enemy.Alive() {
				return fmt.Sprintf(
					"You attack %s for %d damage. %s has %d/%d health left.",
					game.Enemy.Name(),
					damage,
					game.Enemy.Name(),
					game.Enemy.Health(),
					game.Enemy.MaxHealth(),
				), true
			}

			return fmt.Sprintf(
				"You attack %s for %d damage. %s is defeated.",
				game.Enemy.Name(),
				damage,
				game.Enemy.Name(),
			), true
		},
		Description: []string{"Attack the current enemy."},
	})

	// status shows the current state of the enemy in the battle.
	db.RegisterCommand(core.Command[Game]{
		Name: "status",
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
	})
}

func registerDebugCommands(db *core.CommandDB[Game]) {
}
