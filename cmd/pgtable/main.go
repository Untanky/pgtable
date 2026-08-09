package main

import (
	"fmt"
	"os"
	"runtime/debug"

	tea "charm.land/bubbletea/v2"
	flag "github.com/spf13/pflag"
	"github.com/untanky/pgtable/internal/parser"
	"github.com/untanky/pgtable/internal/tui"
	"github.com/untanky/pgtable/render"
)

func main() {
	var (
		showVersion = flag.BoolP("version", "v", false, "Display the application's version")
	)
	flag.Parse()

	if *showVersion {
		displayVersions()
	}

	checkStdin()

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

func checkStdin() {
	info, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error detecting stdin: %v\n", err)
		os.Exit(1)
		return
	}

	if info.Mode()&os.ModeCharDevice != 0 {
		fmt.Fprintln(os.Stderr, "Please pipe in data; interactive mode not supported!")
		os.Exit(1)
		return
	}
}
