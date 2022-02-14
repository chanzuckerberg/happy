package util

import "testing"

func TestTablePrinter(t *testing.T) {
	printer := NewTablePrinter([]string{"foo"})
	printer.AddRow([]string{"bar"})
	printer.Print()
}
