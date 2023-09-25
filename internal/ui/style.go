package ui

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

var (
	tableOptions = table.Options{
		DoNotColorBordersAndSeparators: false,
		DrawBorder:                     true,
		SeparateColumns:                true,
		SeparateFooter:                 false,
		SeparateHeader:                 true,
		SeparateRows:                   true,
	}

	titleOptions = table.TitleOptions{
		Align: text.AlignCenter,
	}

	TouchyStyle = table.Style{
		Name:    "TouchyStyle",
		Box:     table.StyleBoxRounded,
		Color:   table.ColorOptionsDefault,
		Format:  table.FormatOptionsDefault,
		HTML:    table.DefaultHTMLOptions,
		Options: tableOptions,
		Title:   titleOptions,
	}
)
