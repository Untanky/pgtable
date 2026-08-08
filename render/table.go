package render

import (
	"fmt"
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/untanky/pgtable"
)

type Cursor struct {
	row    int
	column int
}

type screen struct {
	width  int
	height int
}

type offset struct {
	horizontal int
	vertical   int
}

type Theme struct {
	Border color.Color
	Null   color.Color

	HighlightBackground color.Color
	HighlightForeground color.Color

	BottomBarBackground color.Color
	BottomBarForeground color.Color
}

type Model struct {
	table  pgtable.Table
	cursor *Cursor
	screen *screen
	offset *offset
	theme  *Theme
}

func NewModel(table pgtable.Table) *Model {
	return new(Model{
		table:  table,
		cursor: new(Cursor),
		screen: new(screen),
		offset: new(offset),
		theme: new(Theme{
			Border: lipgloss.Color("#6e738d"),

			HighlightBackground: lipgloss.Color("#f5a97f"),
			HighlightForeground: lipgloss.Color("#24273a"),

			BottomBarBackground: lipgloss.Color("#494d64"),
			BottomBarForeground: lipgloss.Color("#a5adcb"),
		}),
	})
}

func (model *Model) Render() string {
	rows := model.renderRows().Y(3)

	tableLayers := lipgloss.NewLayer("").AddLayers(
		model.renderBorder('╭', '┬', '╮'),
		model.renderHeader().Y(1),
		model.renderBorder('├', '┼', '┤').Y(2),
		rows,
		model.renderBorder('╰', '┴', '╯').Y(rows.Height()+3),
	)

	comp := lipgloss.NewCompositor(
		tableLayers.X(-model.offset.horizontal).Y(-model.offset.vertical),
		model.renderBottomBar().Y(model.screen.height-1),
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

	borderStyle := lipgloss.NewStyle().Foreground(model.theme.Border)

	return lipgloss.NewLayer(borderStyle.Render(builder.String()))
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
			BorderForeground(model.theme.Border).
			BorderLeft(idx > 0).
			Width(width).
			Align(lipgloss.Center)

		builder.WriteString(cellStyle.Render(column.Name))
	}

	rowStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true).
		BorderForeground(model.theme.Border)

	return lipgloss.NewLayer(rowStyle.Render(builder.String()))
}

func (model *Model) renderRows() *lipgloss.Layer {
	compositor := lipgloss.NewLayer("")
	rowCount := 0

	for i := 0; i < len(model.table.Cells); i += len(model.table.Columns) {
		row := model.table.Cells[i : i+len(model.table.Columns)]
		compositor.AddLayers(model.renderRow(rowCount, row).Y(rowCount))
		rowCount++
	}

	return compositor
}

func (model *Model) renderRow(rowIdx int, row []pgtable.Cell) *lipgloss.Layer {
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
			BorderForeground(model.theme.Border).
			BorderLeft(idx > 0).
			Width(width).
			Padding(0, 1).
			Align(lipgloss.Left)

		if cell.IsNull {
			cellStyle = cellStyle.Faint(true)
			value = "null"
		}

		if model.cursor.row == rowIdx && model.cursor.column == idx {
			cellStyle = cellStyle.
				Background(model.theme.HighlightBackground).
				Foreground(model.theme.HighlightForeground)
		}

		builder.WriteString(cellStyle.Render(value))
	}

	rowStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true).
		BorderForeground(model.theme.Border)

	return lipgloss.NewLayer(rowStyle.Render(builder.String()))
}

func (model *Model) renderBottomBar() *lipgloss.Layer {
	bottomBarStyle := lipgloss.NewStyle().
		Foreground(model.theme.BottomBarForeground).
		Background(model.theme.BottomBarBackground).
		Width(model.screen.width)

	bottomContent := fmt.Sprintf("row %d/%d, col %d/%d (%s)",
		model.cursor.row+1,
		model.table.RowsCount(),
		model.cursor.column+1,
		model.table.ColumnsCount(),
		model.table.Columns[model.cursor.column].Name,
	)

	return lipgloss.NewLayer(bottomBarStyle.Render(bottomContent))
}

func (model *Model) Move(vertical, horizontal int) {
	nextRow := model.cursor.row + vertical
	nextRow = max(0, min(model.table.RowsCount()-1, nextRow))

	nextColumn := model.cursor.column + horizontal
	nextColumn = max(0, min(model.table.ColumnsCount()-1, nextColumn))

	model.cursor.row = nextRow
	model.cursor.column = nextColumn

	model.adjustVerticalOffset()
	model.adjustHorizontalOffset()
}

func (model *Model) ResizeScreen(width, height int) {
	model.screen.width = width
	model.screen.height = height

	model.adjustVerticalOffset()
	model.adjustHorizontalOffset()
}

func (model *Model) adjustVerticalOffset() {
	screenPosition := model.cursor.row - model.offset.vertical
	staticOffset := 5

	if screenPosition+staticOffset > model.screen.height {
		if model.cursor.row+1 == model.table.RowsCount() {
			staticOffset++
		}

		model.offset.vertical = model.cursor.row - model.screen.height + staticOffset
	}

	staticOffset = 3
	if screenPosition+staticOffset < 0 {
		if model.cursor.row == 0 {
			staticOffset = 0
		}
		model.offset.vertical = model.cursor.row + staticOffset
	}
}

func (model *Model) adjustHorizontalOffset() {
	var cursorLeft, cursorRight int

	for idx, column := range model.table.Columns {
		if idx == model.cursor.column {
			cursorRight = cursorLeft + column.Width + 2
			break
		}

		cursorLeft += column.Width + 1
	}

	if cursorRight-cursorLeft > model.screen.width {
		model.offset.horizontal = cursorLeft
		return
	}

	if cursorRight-model.offset.horizontal > model.screen.width {
		model.offset.horizontal = cursorRight - model.screen.width
		return
	}

	if cursorLeft-model.offset.horizontal < 0 {
		model.offset.horizontal = cursorLeft
		return
	}

	totalWidth := 1

	for _, column := range model.table.Columns {
		totalWidth += column.Width + 1
	}

	if totalWidth <= model.screen.width {
		model.offset.horizontal = 0
		return
	}

	if model.offset.horizontal+model.screen.width > totalWidth {
		model.offset.horizontal = totalWidth - model.screen.width
	}
}
