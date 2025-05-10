package tables

import (
	"fmt"
	"io"
	"text/tabwriter"
)

type Model struct {
	Out      io.Writer
	MinWidth int
	TabWidth int
	Padding  int
	PadChar  byte
	Flags    uint
}

func NewModel(out io.Writer, minWidth, tabWidth, padding int, padChar byte, flags uint) *Model {
	return &Model{
		Out:      out,
		MinWidth: minWidth,
		TabWidth: tabWidth,
		Padding:  padding,
		PadChar:  padChar,
		Flags:    flags,
	}
}

func Print(model *Model, headers []string, data [][]string) {
	w := tabwriter.NewWriter(
		model.Out,
		model.MinWidth,
		model.TabWidth,
		model.Padding,
		model.PadChar,
		model.Flags,
	)

	// Print header
	for _, header := range headers {
		fmt.Fprintf(w, "%s\t", header)
	}
	fmt.Fprint(w, "\n")

	// Print data
	for _, row := range data {
		for _, cell := range row {
			fmt.Fprintf(w, "%s\t", cell)
		}

		fmt.Fprint(w, "\n")
	}

	w.Flush()
}
