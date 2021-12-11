package golanglibs

import (
	"strings"

	"github.com/xuri/excelize/v2"
)

type xlsxStruct struct {
	f *excelize.File
}

func getXlsx(path string) *xlsxStruct {
	if !pathExists(path) {
		f := excelize.NewFile()
		err := f.SaveAs(path)
		panicerr(err)
	}
	k, err := excelize.OpenFile(path)
	panicerr(err)
	return &xlsxStruct{f: k}
}

type xlsxSheetStruct struct {
	x     *xlsxStruct
	sheet string
}

func (c *xlsxStruct) getSheet(name string) *xlsxSheetStruct {
	name = strings.Title(name)
	if !Array(c.f.GetSheetList()).Has(name) {
		Lg.trace("make new sheet")
		c.f.NewSheet(name)
	}
	return &xlsxSheetStruct{
		x:     c,
		sheet: name,
	}
}

func (c *xlsxSheetStruct) get(coordinate string) string {
	s, err := c.x.f.GetCellValue(c.sheet, coordinate)
	panicerr(err)
	return s
}

func (c *xlsxSheetStruct) set(coordinate string, value string) *xlsxSheetStruct {
	err := c.x.f.SetCellValue(c.sheet, coordinate, value)
	panicerr(err)
	return c
}

func (c *xlsxStruct) save() {
	c.f.Save()
}

func (c *xlsxStruct) close() {
	c.save()
}
