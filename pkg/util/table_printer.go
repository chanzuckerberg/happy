package util

import "fmt"

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

type TablePrinter struct {
	rows     [][]string
	widths   []int
	headings []string
}

func NewTablePrinter(headings []string) *TablePrinter {
	return &TablePrinter{
		headings: headings,
	}
}

func (s *TablePrinter) BumpWidth(data []string) {
	for i, entry := range data {
		if i >= len(s.widths) {
			s.widths = append(s.widths, len(data[i]))
		} else {
			s.widths[i] = Max(len(entry), s.widths[i])
		}
	}
}

func (s *TablePrinter) AddRow(data []string) {
	s.BumpWidth(data)
	s.rows = append(s.rows, data)
}

func (s *TablePrinter) Print() {
	var fmtString string
	for _, width := range s.widths {
		fmtString += fmt.Sprintf("%%%dv  ", -width)
	}
	fmtString += "\n"

	headings := make([]interface{}, len(s.headings))
	for i, v := range s.headings {
		headings[i] = v
	}
	fmt.Printf(fmtString, headings...)

	separators := make([]interface{}, len(s.headings))
	for i := range separators {
		separators[i] = "-----"
	}
	fmt.Printf(fmtString, separators...)

	for _, row := range s.rows {
		iRow := make([]interface{}, len(row))
		for i, v := range row {
			iRow[i] = v
		}
		fmt.Printf(fmtString, iRow...)
	}
}
