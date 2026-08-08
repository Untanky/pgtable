package tui

import (
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (app App) handleKeys(msg tea.KeyMsg) inputState {
	key := msg.String()
	state := app.state

	if len(key) == 1 && key >= "0" && key <= "9" {
		if key != "0" || state.count != "" {
			state.count += key
			return state
		}
	}

	state.pending += key
	seq := state.pending

	if cmd, ok := app.commands[seq]; ok {
		n := state.getCount()
		app.execute(cmd, n)
		return inputState{}
	}

	// Partial match -> keep waiting
	for prefix := range app.commands {
		if strings.HasPrefix(prefix, seq) {
			return state
		}
	}

	// no match; abort sequence
	return inputState{}
}

func (app App) execute(cmd string, n int) {
	switch cmd {
	case "left":
		app.tableModel.Move(0, -n)
	case "rigth":
		app.tableModel.Move(0, n)
	case "up":
		app.tableModel.Move(-n, 0)
	case "down":
		app.tableModel.Move(n,0)
	case "yank":
		app.tableModel.Yank()
	}
}

func (state inputState) getCount() int {
	if state.count == "" {
		return 1
	}
	n, _ := strconv.Atoi(state.count)
	return n
}
