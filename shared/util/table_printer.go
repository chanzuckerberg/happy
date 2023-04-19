package util

import (
	"os"
	"sync"

	"github.com/lensesio/tableprinter"
)

type row struct {
	Resource string `header:"Resource"`
	Value    string `header:"Value"`
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

type TablePrinter struct {
	once    *sync.Once
	printer *tableprinter.Printer
	rows    []interface{}
}

func NewTablePrinter() *TablePrinter {
	printer := tableprinter.New(os.Stderr)
	printer.RowCharLimit = 60
	printer.AutoWrapText = true
	printer.BorderTop, printer.BorderBottom, printer.BorderLeft, printer.BorderRight = true, true, true, true
	printer.CenterSeparator = "│"
	printer.ColumnSeparator = "│"
	printer.RowSeparator = "─"
	return &TablePrinter{
		once:    &sync.Once{},
		printer: printer,
	}
}

func row2Console(resouce string, value string) row {
	return row{Resource: resouce, Value: value}
}

// Adds a structured row to the cache
func (s *TablePrinter) AddRow(row interface{}) {
	s.rows = append(s.rows, row)
}

// Adds a two column row to the cache
func (s *TablePrinter) AddSimpleRow(resource string, value string) {
	s.rows = append(s.rows, row2Console(resource, value))
}

// Directly prints the rows passed in
func (s *TablePrinter) Print(in interface{}) {
	s.printer.Print(in)
}

// Flushes out the row cache, printing the them all out
func (s *TablePrinter) Flush() {
	s.once.Do(func() { s.Print(s.rows) })
}
