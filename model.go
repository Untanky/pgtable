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
