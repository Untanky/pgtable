package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/untanky/pgtable/internal/parser"
	"github.com/untanky/pgtable/internal/tui"
	"github.com/untanky/pgtable/render"
)

func main() {
	table, err := parser.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	renderer := render.NewModel(table)
	program := tea.NewProgram(tui.NewApp(*renderer))

	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
