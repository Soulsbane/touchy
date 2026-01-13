package ui

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func CreateNewTableWriter(title string, args ...any) table.Writer {
	writer := table.NewWriter()

	writer.SetOutputMirror(os.Stdout)
	writer.SetTitle(title)
	writer.AppendHeader(args)
	writer.SetStyle(TouchyStyle)
	writer.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignLeft, WidthMin: 30},
	})

	return writer
}
