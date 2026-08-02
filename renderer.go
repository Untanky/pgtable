package pgtable

import (
	"bytes"
	"io"
	"strings"
)

func Render(writer io.Writer, table Table) error {
	totalWidth := len(table.Columns) + 1
	for _, column := range table.Columns {
		totalWidth += column.Width
	}

	firstLine := bytes.NewBuffer(nil)
	header := bytes.NewBuffer(nil)
	separator := bytes.NewBuffer(nil)
	content := bytes.NewBuffer(nil)
	lastLine := bytes.NewBuffer(nil)

	for idx, column := range table.Columns {
		if idx == 0 {
			firstLine.WriteRune('╭')
			separator.WriteRune('├')
			lastLine.WriteRune('╰')
		} else {
			firstLine.WriteRune('┬')
			separator.WriteRune('┼')
			lastLine.WriteRune('┴')
		}

		header.WriteRune('│')

		header.WriteString(padCenter(column.Name, column.Width, " "))
		firstLine.Write(repeatRune('─', column.Width))
		separator.Write(repeatRune('─', column.Width))
		lastLine.Write(repeatRune('─', column.Width))

		if idx == len(table.Columns)-1 {
			firstLine.WriteRune('╮')
			header.WriteRune('│')
			separator.WriteRune('┤')
			lastLine.WriteRune('╯')
		}
	}

	for idx, column := range table.Cells {
		columnIdx := idx % len(table.Columns)

		content.WriteString("│ ")
		content.WriteString(padRight(column.Value, table.Columns[columnIdx].Width-1, " "))

		if columnIdx == len(table.Columns)-1 {
			content.WriteString("│\n")
		}
	}

	firstLine.WriteByte('\n')
	header.WriteByte('\n')
	separator.WriteByte('\n')
	lastLine.WriteByte('\n')

	firstLine.WriteTo(writer)
	header.WriteTo(writer)
	separator.WriteTo(writer)
	content.WriteTo(writer)
	lastLine.WriteTo(writer)

	return nil
}

func repeatRune(r rune, count int) []byte {
	return bytes.Repeat([]byte(string(r)), count)
}

func padRight(s string, length int, pad string) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(pad, length-len(s))
}

func padCenter(s string, length int, pad string) string {
	if len(s) >= length {
		return s
	}
	total := length - len(s)
	left := total / 2
	right := total - left
	return strings.Repeat(pad, left) + s + strings.Repeat(pad, right)
}

