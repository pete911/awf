package out

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

const maxColumnSize = 25

type Table struct {
	writer *tabwriter.Writer
	trim   bool
}

func NewTable(output io.Writer, trim bool) Table {
	return Table{
		writer: tabwriter.NewWriter(output, 1, 1, 2, ' ', 0),
		trim:   trim,
	}
}

func (t Table) AddRow(columns ...string) {
	if t.trim {
		for i := range columns {
			columns[i] = squashTo(columns[i], maxColumnSize)
		}
	}
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

func trimTo(in string, max int) string {
	if len(in) < max {
		return in
	}
	return fmt.Sprintf("%s..", in[:max-2])
}

func squashTo(in string, max int) string {
	if len(in) < max {
		return in
	}

	leftMax := max / 2
	if max%2 == 0 {
		leftMax--
	}
	rightMax := len(in) - (max / 2) + 1
	return fmt.Sprintf("%s..%s", in[:leftMax], in[rightMax:])
}
