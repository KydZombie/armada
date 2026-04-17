package game

import (
	"fmt"

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
	})
}

func registerDebugCommands(db *core.CommandDB[Game]) {
}
