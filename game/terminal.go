package game

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/KydZombie/armada/core"
	"github.com/KydZombie/armada/core/util"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type historyType int8

const (
	userInputType historyType = iota
	standardReturnType
	errorReturnType
)

type terminalHistoryItem struct {
	historyType historyType
	text        string
}

type TerminalWindow struct {
	core.BaseWindow[Game]

	commandDB *core.CommandDB[Game]

	inputText strings.Builder
	history   []terminalHistoryItem
}

func NewTerminalWindow(sizeFunc func(gm *core.GameManager) rl.Rectangle, gm *core.GameManager, commandDB *core.CommandDB[Game]) *TerminalWindow {
	return &TerminalWindow{
		BaseWindow: core.NewBaseWindow[Game](sizeFunc, gm, false),

		commandDB: commandDB,

		inputText: strings.Builder{},
	}
}

func (t *TerminalWindow) HandleInput(gm *core.GameManager, state *Game) bool {
	captured := t.IsVisible()

	if rl.IsKeyPressed(rl.KeyTab) {
		t.SetVisible(!t.IsVisible())
		captured = true
	}

	if t.IsVisible() {
		key := rl.GetKeyPressed()
		if key == rl.KeyEnter {
			cmd := t.inputText.String()
			t.handleCommand(gm, state, cmd)
			t.inputText.Reset()
		} else if key == rl.KeyBackspace {
			if len(t.inputText.String()) > 0 {
				curr := t.inputText.String()
				t.inputText.Reset()
				_, size := utf8.DecodeLastRuneInString(curr)
				t.inputText.WriteString(curr[:len(curr)-size])
			}
		} else {
			r, isRune := util.KeyToAlphanumeric(key)
			if isRune {
				t.inputText.WriteRune(r)
			}
		}

		key = rl.GetKeyPressed()
	}

	return captured
}

func (t *TerminalWindow) UpdateWindow(gm *core.GameManager, state *Game) {

}

func (t *TerminalWindow) DrawWindow(gm *core.GameManager, state *Game) {
	if !t.IsVisible() {
		return
	}

	bounds := t.GetBounds()
	offsetX := int32(bounds.X) + 4
	offsetY := int32(bounds.Y) + 4

	rl.DrawRectangleRec(bounds, rl.Black)
	rl.DrawText(fmt.Sprint("> ", t.inputText.String()), offsetX, offsetY, 24, rl.White)

	for i := range t.history {
		historyItem := t.history[len(t.history)-i-1]
		var color rl.Color
		switch historyItem.historyType {
		case userInputType:
			color = rl.DarkGray
		case standardReturnType:
			color = rl.LightGray
		case errorReturnType:
			color = rl.Red
		}
		rl.DrawText(historyItem.text, offsetX, offsetY+int32(i+1)*24, 24, color)
	}
}

func (t *TerminalWindow) DrawWindowUI(gm *core.GameManager, state *Game) {

}

func (t *TerminalWindow) handleCommand(gm *core.GameManager, state *Game, rawCmd string) {
	gm.Log.Println("Command: ", rawCmd)

	t.history = append(t.history, terminalHistoryItem{
		historyType: userInputType,
		text:        fmt.Sprint("> ", t.inputText.String()),
	})

	result, success := t.commandDB.ParseAndRunCommand(rawCmd, state)
	if len(result) != 0 {
		var historyType historyType
		if success {
			historyType = standardReturnType
		} else {
			historyType = errorReturnType
		}

		t.history = append(t.history, terminalHistoryItem{
			historyType: historyType,
			text:        result,
		})
	}
}
