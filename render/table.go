package render

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/untanky/pgtable"
)

type Cursor struct {
	row    int
	column int
}

type Model struct {
	table  pgtable.Table
	cursor *Cursor
}

func NewModel(table pgtable.Table) *Model {
	return new(Model{
		table:  table,
		cursor: new(Cursor),
	})
}

func (model *Model) Render() string {

	comp := lipgloss.NewCompositor(
		model.renderBorder('╭', '┬', '╮'),
		model.renderHeader().Y(1),
		model.renderBorder('├', '┼', '┤').Y(2),
		model.renderRows().Y(3),
		model.renderBorder('╰', '┴', '╯').Y(4),
	)
	return comp.Render()
}

func (model *Model) GetCursor() *Cursor {
	return model.cursor
}

func (model *Model) renderBorder(left, middle, right rune) *lipgloss.Layer {
	var builder strings.Builder

	for idx, column := range model.table.Columns {
		if idx == 0 {
			builder.WriteRune(left)
		} else {
			builder.WriteRune(middle)
		}

		builder.WriteString(strings.Repeat("─", column.Width))

		if idx == len(model.table.Columns)-1 {
			builder.WriteRune(right)
		}
	}

	return lipgloss.NewLayer(builder.String())
}

func (model *Model) renderHeader() *lipgloss.Layer {
	var builder strings.Builder

	for idx, column := range model.table.Columns {
		width := column.Width + 1
		if idx == 0 {
			width = column.Width
		}

		cellStyle := lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderLeft(idx > 0).
			Width(width).
			Align(lipgloss.Center)

		builder.WriteString(cellStyle.Render(column.Name))
	}

	rowStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, true)

	return lipgloss.NewLayer(rowStyle.Render(builder.String()))
}

func (model *Model) renderRows() *lipgloss.Layer {
	layers := []*lipgloss.Layer{}

	for i := 0; i < len(model.table.Cells); i += len(model.table.Columns) {
		layers = append(layers, model.renderRow(model.table.Cells[i:i+len(model.table.Columns)]))
	}

	return layers[0]
}

func (model *Model) renderRow(row []pgtable.Cell) *lipgloss.Layer {
	var builder strings.Builder

	for idx, cell := range row {
		column := model.table.Columns[idx]

		width := column.Width + 1
		if idx == 0 {
			width = column.Width
		}

		value := cell.Value

		cellStyle := lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderLeft(idx > 0).
			Width(width).
			Align(lipgloss.Center)

		if cell.IsNull {
			cellStyle = cellStyle.Foreground(lipgloss.Color("249"))
			value = "null"
		}

		builder.WriteString(cellStyle.Render(value))
	}

	rowStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, true)

	return lipgloss.NewLayer(rowStyle.Render(builder.String()))
}
