package core

import (
	"fmt"
	"strings"
)

type Command[State any] struct {
	Name  string
	Info  string
	OnRun func(args []string, state *State) (string, bool)
}

type CommandDB[State any] struct {
	cmdMap map[string]Command[State]
}

func NewCommandDB[State any]() *CommandDB[State] {
	return &CommandDB[State]{
		cmdMap: make(map[string]Command[State]),
	}
}

func (db *CommandDB[State]) getCommands() map[string]Command[State] {
	return db.cmdMap
}

func (db *CommandDB[State]) RegisterCommand(command Command[State]) {
	db.cmdMap[command.Name] = command
}

// ParseAndRunCommand returns the result of the command as (output, success)
func (db *CommandDB[State]) ParseAndRunCommand(fullCommand string, state *State) (string, bool) {
	args := strings.Split(fullCommand, " ")
	cmdName := args[0]
	args = args[1:]

	cmd, ok := db.cmdMap[cmdName]
	if !ok {
		return fmt.Sprint("Unknown command, ", cmdName), false
	}

	return cmd.OnRun(args, state)
}
