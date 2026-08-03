package pgtable

type Cell struct {
	Value  string
	IsNull bool
}

type Column struct {
	Name     string
	Width    int // display width in terminal cells
	IsNumber bool
}

type Table struct {
	Columns []Column
	Cells   []Cell
	RowHint int // parsed from "(N rows)", 0 if absent
}

func (t Table) ColumnsCount() int {
	return len(t.Columns)
}

func (t Table) RowsCount() int {
	return len(t.Cells) / len(t.Columns)
}
