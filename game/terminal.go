package game

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/KydZombie/armada/core"
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
	if state.isGameOverModalActive() || state.isMissionBriefingActive() {
		return false
	}

	for key := rl.GetKeyPressed(); key != 0; key = rl.GetKeyPressed() {
		switch key {
		case rl.KeyUp:
			if t.currentHistoryLine < int32(len(t.history)) {
				t.currentHistoryLine++
			}
		case rl.KeyDown:
			if t.currentHistoryLine > 0 {
				t.currentHistoryLine--
			}
		case rl.KeyEnter:
			cmd := t.inputText.String()
			t.handleCommand(gm, state, cmd)
			t.inputText.Reset()
		case rl.KeyBackspace:
			if len(t.inputText.String()) > 0 {
				curr := t.inputText.String()
				t.inputText.Reset()
				_, size := utf8.DecodeLastRuneInString(curr)
				t.inputText.WriteString(curr[:len(curr)-size])
			}
		}
	}

	for char := rl.GetCharPressed(); char != 0; char = rl.GetCharPressed() {
		switch {
		case char >= 'A' && char <= 'Z':
			t.inputText.WriteRune(char + ('a' - 'A'))
		case char >= 32 && char <= 126:
			t.inputText.WriteRune(char)
		}
	}

	return true
}

func (t *TerminalWindow) UpdateWindow(gm *core.GameManager, state *Game) {

}

func wrapTerminalLine(text string, maxWidth float32, fontSize int32) []string {
	if text == "" {
		return []string{""}
	}

	if maxWidth <= 0 {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{text}
	}

	lines := make([]string, 0, len(words))
	currentLine := ""

	appendChunkedWord := func(word string) {
		runes := []rune(word)
		start := 0
		for start < len(runes) {
			end := start + 1
			for end <= len(runes) {
				chunk := string(runes[start:end])
				if float32(rl.MeasureText(chunk, fontSize)) > maxWidth {
					end--
					break
				}
				end++
			}

			if end <= start {
				end = start + 1
			}

			lines = append(lines, string(runes[start:end]))
			start = end
		}
	}

	for _, word := range words {
		candidate := word
		if currentLine != "" {
			candidate = currentLine + " " + word
		}

		if float32(rl.MeasureText(candidate, fontSize)) <= maxWidth {
			currentLine = candidate
			continue
		}

		if currentLine != "" {
			lines = append(lines, currentLine)
			currentLine = ""
		}

		if float32(rl.MeasureText(word, fontSize)) <= maxWidth {
			currentLine = word
			continue
		}

		appendChunkedWord(word)
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	if len(lines) == 0 {
		return []string{text}
	}

	return lines
}

func (t *TerminalWindow) DrawWindow(gm *core.GameManager, state *Game) {
	if !t.IsVisible() {
		return
	}
	const fontSize int32 = 18

	const innerOffset int32 = 4

	bounds := t.GetBounds()
	offsetX := int32(bounds.X) + innerOffset
	offsetY := int32(bounds.Y) + innerOffset
	contentWidth := bounds.Width - float32(innerOffset*2)

	rl.DrawRectangleRec(bounds, rl.Black)
	rl.BeginScissorMode(int32(bounds.X), int32(bounds.Y), int32(bounds.Width), int32(bounds.Height))

	inputLines := wrapTerminalLine(fmt.Sprint("> ", t.inputText.String()), contentWidth, fontSize)
	for i, line := range inputLines {
		rl.DrawText(line, offsetX, offsetY+int32(i)*fontSize, fontSize, rl.White)
	}

	maxVisibleHeight := int32(bounds.Height) - innerOffset
	nextY := offsetY + int32(len(inputLines))*fontSize
	for lineNumber := t.currentHistoryLine - 1; lineNumber >= 0 && nextY+fontSize <= maxVisibleHeight+int32(bounds.Y); lineNumber-- {
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

		wrappedLines := wrapTerminalLine(historyItem.text, contentWidth, fontSize)
		for _, line := range wrappedLines {
			if nextY+fontSize > maxVisibleHeight+int32(bounds.Y) {
				break
			}
			rl.DrawText(line, offsetX, nextY, fontSize, color)
			nextY += fontSize
		}
	}

	rl.EndScissorMode()
}
