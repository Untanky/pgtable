package render

import "github.com/atotto/clipboard"

func (model *Model) Yank() error {
	rowPrefix := model.cursor.row * model.table.ColumnsCount()
	cellIndex := rowPrefix + model.cursor.column

	value := model.table.Cells[cellIndex]

	if value.IsNull {
		return nil
	}

	return clipboard.WriteAll(value.Value)
}
