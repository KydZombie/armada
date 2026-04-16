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

}

func registerCombatCommands(db *core.CommandDB[Game]) {

}

func registerDebugCommands(db *core.CommandDB[Game]) {
}
