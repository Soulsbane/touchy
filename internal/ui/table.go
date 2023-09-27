package ui

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
)

func CreateNewTableWriter(title string, args ...interface{}) table.Writer {
	writer := table.NewWriter()

	writer.SetOutputMirror(os.Stdout)
	writer.SetTitle(title)
	writer.AppendHeader(args)
	writer.SetStyle(TouchyStyle)

	return writer
}
