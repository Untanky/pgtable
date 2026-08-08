package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/untanky/pgtable/render"
)

type inputState struct {
	count   string
	pending string
}

type App struct {
	state      inputState
	commands   map[string]string
	tableModel render.Model
}

func NewApp(tableModel render.Model) *App {
	return new(App{
		state: inputState{},
		commands: map[string]string{
			"ctrl+c": "quit",
			"q":      "quit",
			"h":      "left",
			"left":   "left",
			"l":      "right",
			"right":  "right",
			"j":      "down",
			"down":   "down",
			"k":      "up",
			"up":     "up",
			"yy":     "yank",
			"ctrl+d": "down-half-screen",
			"ctrl+u": "up-half-screen",
		},
		tableModel: tableModel,
	})
}

func (app App) Init() tea.Cmd {
	return nil
}

func (app App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		app.tableModel.ResizeScreen(msg.Width, msg.Height)

	// Is it a key press?
	case tea.KeyPressMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return app, tea.Quit

		default:
			app.state = app.handleKeys(msg)
			return app, nil
		}
	}

	return app, nil
}

func (app App) View() tea.View {
	view := tea.NewView(app.tableModel.Render())
	view.AltScreen = true

	return view
}
