package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTablePrinter(t *testing.T) {
	r := require.New(t)
	printer := NewTablePrinter()
	printer.Print([]string{"bar"})
	r.Equal(5, Max(2, 5))
	r.Equal(5, Max(5, 2))
}
