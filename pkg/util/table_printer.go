package util

import (
	"strings"
	"sync"

	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
)

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

type TablePrinter struct {
	once *sync.Once

	table  *tablewriter.Table
	buffer *strings.Builder
}

func NewTablePrinter(headings []string) *TablePrinter {
	buffer := &strings.Builder{}
	table := tablewriter.NewWriter(buffer)

	table.SetHeader(headings)
	return &TablePrinter{
		once: &sync.Once{},

		table:  table,
		buffer: buffer,
	}
}

func (s *TablePrinter) AddRow(data []string) {
	s.table.Append(data)
}

func (s *TablePrinter) Print() {
	s.once.Do(func() { s.table.Render() })
	logrus.Printf("\n %s", s.buffer.String())
}
