package golanglibs

import (
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

type progressBarStruct struct {
	bar *progressbar.ProgressBar
}

func getProgressBar(title string, total int64, showBytes ...bool) *progressBarStruct {
	var showBytesOption progressbar.Option
	if len(showBytes) != 0 && showBytes[0] == true {
		showBytesOption = progressbar.OptionShowBytes(true)
	} else {
		showBytesOption = progressbar.OptionShowBytes(false)
	}
	bar := progressbar.NewOptions64(total,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionFullWidth(),
		showBytesOption,
		progressbar.OptionSetDescription("[cyan]*[reset] "+title),

		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[cyan]=[reset]",
			SaucerHead:    "[cyan]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	return &progressBarStruct{
		bar: bar,
	}
}

func (m *progressBarStruct) Add(num int64) {
	m.bar.Add64(num)
}

func (m *progressBarStruct) Set(num int64) {
	m.bar.Set64(num)
}

func (m *progressBarStruct) SetTotal(total int64) {
	m.bar.ChangeMax64(total)
}

func (m *progressBarStruct) Clear() {
	m.bar.Clear()
}
