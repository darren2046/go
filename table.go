package golanglibs

import (
	"github.com/alexeyco/simpletable"
	"github.com/bndr/gotabulate"
)

type tableStruct struct {
	row          [][]interface{}
	header       []string
	maxCellWidth int
}

func getTable(header ...string) *tableStruct {
	return &tableStruct{header: header, maxCellWidth: 0}
}

func (m *tableStruct) setMaxCellWidth(width ...int) {
	if len(width) == 0 {
		m.maxCellWidth = 30
	} else {
		m.maxCellWidth = width[0]
	}
}

func (m *tableStruct) addRow(row ...interface{}) {
	if len(row) != len(m.header) {
		panicerr("添加的数据个数跟表头的个数对不上")
	}
	m.row = append(m.row, row)
}

func (m *tableStruct) render() string {
	if m.maxCellWidth == 0 {
		table := simpletable.New()

		for _, header := range m.header {
			table.Header.Cells = append(table.Header.Cells, &simpletable.Cell{Align: simpletable.AlignLeft, Text: header})
		}

		for _, row := range m.row {
			var cell []*simpletable.Cell
			for _, r := range row {
				cell = append(cell, &simpletable.Cell{
					Align: simpletable.AlignLeft,
					Text:  Str(r),
				})
			}

			table.Body.Cells = append(table.Body.Cells, cell)
		}

		table.SetStyle(simpletable.StyleCompactLite)
		return table.String()
	}

	tabulate := gotabulate.Create(m.row)
	tabulate.SetHeaders(m.header)
	tabulate.SetWrapStrings(true)
	tabulate.SetMaxCellSize(m.maxCellWidth)
	return tabulate.Render("grid")
}
