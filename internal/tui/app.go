package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/untanky/pgtable/render"
)

type App struct {
	tableModel render.Model
}

func NewApp(tableModel render.Model) *App {
	return new(App{
		tableModel: tableModel,
	})
}

func (app *App) Init() tea.Cmd {
	return nil
}

func (app *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyPressMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return app, tea.Quit
		}
	}

	return app, nil
}

func (app *App) View() tea.View {
	view := tea.NewView(app.tableModel.Render())
	view.AltScreen = true

	return view
}
