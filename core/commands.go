package core

import (
	"fmt"
	"strings"
)

type Command[State any] struct {
	Name        string
	Info        string
	OnRun       func(args []string, state *State) (string, bool)
	Description []string
}

type CommandDB[State any] struct {
	CmdMap map[string]Command[State]
}

func NewCommandDB[State any]() *CommandDB[State] {
	return &CommandDB[State]{
		CmdMap: make(map[string]Command[State]),
	}
}

func (db *CommandDB[State]) getCommands() map[string]Command[State] {
	return db.CmdMap
}

func (db *CommandDB[State]) RegisterCommand(command Command[State]) {
	db.CmdMap[command.Name] = command
}

// ParseAndRunCommand returns the result of the command as (output, success)
func (db *CommandDB[State]) ParseAndRunCommand(fullCommand string, state *State) (string, bool) {
	args := strings.Fields(fullCommand)
	if len(args) == 0 {
		return "", true
	}

	cmdName := strings.ToLower(args[0])
	args = args[1:]

	cmd, ok := db.CmdMap[cmdName]
	if !ok {
		return fmt.Sprint("Unknown command, ", cmdName), false
	}

	return cmd.OnRun(args, state)
}
