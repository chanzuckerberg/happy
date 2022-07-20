package util

import (
	"os"
	"sync"

	"github.com/lensesio/tableprinter"
)

type Row struct {
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

func Row2Console(resouce string, value string) Row {
	return Row{Resource: resouce, Value: value}
}

func (s *TablePrinter) AddRow(row interface{}) {
	s.rows = append(s.rows, row)
}

func (s *TablePrinter) AddSimpleRow(resource string, value string) {
	s.rows = append(s.rows, Row2Console(resource, value))
}

func (s *TablePrinter) Print(in interface{}) {
	s.printer.Print(in)
}

func (s *TablePrinter) Flush() {
	s.once.Do(func() { s.printer.Print(s.rows) })
}
