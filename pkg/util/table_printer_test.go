package util

import "testing"

func TestTablePrinter(t *testing.T) {
	printer := NewTablePrinter([]string{"foo"})
	printer.AddRow("bar")
	printer.Print()
}
