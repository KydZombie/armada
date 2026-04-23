package game

import (
	"fmt"
	"strings"

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
			return "You can type commands here!", true
		},
	})
}

func registerInfoCommands(db *core.CommandDB[Game]) {

}

func registerUnitCommands(db *core.CommandDB[Game]) {

}

func registerCombatCommands(db *core.CommandDB[Game]) {
	// attack is a very small prototype combat command.
	// It always deals 1 damage to the current enemy so we can test
	// terminal commands and battle window updates together.
	db.RegisterCommand(core.Command[Game]{
		Name: "attack",
		OnRun: func(args []string, game *Game) (string, bool) {
			targetPart := ""
			if len(args) > 0 {
				targetPart = args[0]
			}

			return game.FireSelectedWeapon(targetPart)
		},
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
	})

	// parts shows the enemy's body parts and their health values.
	// This is mainly for debugging now, and it will also help when
	// targeted combat is added later.
	db.RegisterCommand(core.Command[Game]{
		Name: "parts",
		OnRun: func(args []string, game *Game) (string, bool) {
			if game.Enemy == nil {
				return "There is no enemy to inspect.", false
			}

			enemy, ok := game.Enemy.(*BasicEnemy)
			if !ok {
				return "Enemy parts are not available.", false
			}

			var output strings.Builder
			firstPart := true

			for _, part := range enemy.Parts {
				if part == nil {
					continue
				}

				if !firstPart {
					output.WriteString(" ")
				}
				firstPart = false

				output.WriteString(fmt.Sprintf(
					"%s:%d/%d",
					part.Name,
					part.Health,
					part.MaxHealth,
				))
			}

			return output.String(), true
		},
	})

	// target changes which enemy is selected in the multi-enemy prototype.
	db.RegisterCommand(core.Command[Game]{
		Name: "target",
		OnRun: func(args []string, game *Game) (string, bool) {
			if len(args) == 0 {
				return "Choose a target: a, b, or c.", false
			}

			targetLabel := strings.ToLower(args[0])
			targetIndex := -1

			switch targetLabel {
			case "a":
				targetIndex = 0
			case "b":
				targetIndex = 1
			case "c":
				targetIndex = 2
			default:
				return "That target is not valid.", false
			}

			if targetIndex < 0 || targetIndex >= len(game.Enemies) {
				return "That target is out of range.", false
			}

			game.SelectRoom(targetIndex)
			game.SelectionPopupText = fmt.Sprintf("Selected: [%s]", strings.ToUpper(targetLabel))
			game.SelectionPopupFrames = 120

			return fmt.Sprintf("Selected enemy %s.", strings.ToUpper(targetLabel)), true
		},
	})
}

func registerDebugCommands(db *core.CommandDB[Game]) {
}
