package out

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

type Table struct {
	writer *tabwriter.Writer
}

func NewTable(output io.Writer) Table {
	return Table{
		writer: tabwriter.NewWriter(output, 1, 1, 2, ' ', 0),
	}
}

func (t Table) AddRow(columns ...string) {
	if _, err := fmt.Fprintln(t.writer, strings.Join(columns, "\t")); err != nil {
		fmt.Printf("table: add row: %v", err)
	}
}

func (t Table) Print() {
	if err := t.writer.Flush(); err != nil {
		fmt.Printf("table: print: %v", err)
	}
}

func FromInt(in int) string {
	if in == 0 {
		return "-"
	}
	return fmt.Sprintf("%d", in)
}

func TrimTo(in string, max int) string {
	if len(in) < max {
		return in
	}
	return fmt.Sprintf("%s...", in[:max-3])
}
