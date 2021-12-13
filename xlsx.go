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
		Panicerr(err)
	}
	k, err := excelize.OpenFile(path)
	Panicerr(err)
	return &xlsxStruct{f: k}
}

type xlsxSheetStruct struct {
	x     *xlsxStruct
	sheet string
}

func (c *xlsxStruct) GetSheet(name string) *xlsxSheetStruct {
	name = strings.Title(name)
	if !Array(c.f.GetSheetList()).Has(name) {
		Lg.Trace("make new sheet")
		c.f.NewSheet(name)
	}
	return &xlsxSheetStruct{
		x:     c,
		sheet: name,
	}
}

func (c *xlsxSheetStruct) Get(coordinate string) string {
	s, err := c.x.f.GetCellValue(c.sheet, coordinate)
	Panicerr(err)
	return s
}

func (c *xlsxSheetStruct) Set(coordinate string, value string) *xlsxSheetStruct {
	err := c.x.f.SetCellValue(c.sheet, coordinate, value)
	Panicerr(err)
	return c
}

func (c *xlsxStruct) Save() {
	c.f.Save()
}

func (c *xlsxStruct) Close() {
	c.Save()
}
