package game

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/KydZombie/armada/core"
	"github.com/KydZombie/armada/core/util"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type TerminalMessageType int8

const (
	UserInputMessage TerminalMessageType = iota
	StandardOutputMessage
	ErrorMessageMessage
)

type terminalHistoryItem struct {
	historyType TerminalMessageType
	text        string
}

type Terminal struct {
	commandDB *core.CommandDB[Game]

	inputText strings.Builder

	history            []terminalHistoryItem
	currentHistoryLine int32
}

func (t *Terminal) OutputText(text string, messageType TerminalMessageType) {
	t.history = append(t.history, terminalHistoryItem{historyType: messageType, text: text})
	t.currentHistoryLine++
}

func (t *Terminal) handleCommand(gm *core.GameManager, state *Game, rawCmd string) {
	gm.Log.Println("Command: ", rawCmd)

	t.OutputText(fmt.Sprint("> ", t.inputText.String()), UserInputMessage)

	result, success := t.commandDB.ParseAndRunCommand(rawCmd, state)
	if len(result) != 0 {
		var messageType TerminalMessageType
		if success {
			messageType = StandardOutputMessage
		} else {
			messageType = ErrorMessageMessage
		}

		t.OutputText(result, messageType)
	}
}

type TerminalWindow struct {
	core.BaseWindow[Game]
	*Terminal
}

func NewTerminalWindow(sizeFunc func(gm *core.GameManager) rl.Rectangle, gm *core.GameManager, terminal *Terminal) *TerminalWindow {
	return &TerminalWindow{
		BaseWindow: core.NewBaseWindow[Game](sizeFunc, gm, true),
		Terminal:   terminal,
	}
}

func (t *TerminalWindow) HandleInput(gm *core.GameManager, state *Game) bool {
	if !t.IsVisible() {
		return false
	}

	key := rl.GetKeyPressed()
	if key == rl.KeyUp {
		if t.currentHistoryLine < int32(len(t.history)) {
			t.currentHistoryLine++
		}
	} else if key == rl.KeyDown {
		if t.currentHistoryLine > 0 {
			t.currentHistoryLine--
		}
	} else if key == rl.KeyEnter {
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

	return true
}

func (t *TerminalWindow) UpdateWindow(gm *core.GameManager, state *Game) {

}

func (t *TerminalWindow) DrawWindow(gm *core.GameManager, state *Game) {
	if !t.IsVisible() {
		return
	}
	const fontSize int32 = 24

	const innerOffset int32 = 4

	bounds := t.GetBounds()
	offsetX := int32(bounds.X) + innerOffset
	offsetY := int32(bounds.Y) + innerOffset

	rl.DrawRectangleRec(bounds, rl.Black)
	rl.DrawText(fmt.Sprint("> ", t.inputText.String()), offsetX, offsetY, fontSize, rl.White)

	spaceToShowHistory := bounds.Height - float32(innerOffset) - float32(fontSize)
	maxLinesToShow := int32(spaceToShowHistory / float32(fontSize))
	for i := range maxLinesToShow {
		lineNumber := t.currentHistoryLine - 1 - i
		if lineNumber >= int32(len(t.history)) || lineNumber < 0 {
			break
		}

		historyItem := t.history[lineNumber]
		var color rl.Color
		switch historyItem.historyType {
		case UserInputMessage:
			color = rl.DarkGray
		case StandardOutputMessage:
			color = rl.LightGray
		case ErrorMessageMessage:
			color = rl.Red
		}
		rl.DrawText(historyItem.text, offsetX, offsetY+int32(i+1)*24, 24, color)
	}
}

func (t *TerminalWindow) DrawWindowUI(gm *core.GameManager, state *Game) {

}
