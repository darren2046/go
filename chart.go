package golanglibs

import (
	"os"
	"time"

	"github.com/wcharczuk/go-chart"
)

type chartStruct struct {
	LineChartWithTimestampAndNumber func(timestampX []int64, dataY []float64, xtitle string, ytitle string, title string, fpath string)
	LineChartWithNumberAndNumber    func(dataX []float64, dataY []float64, xtitle string, ytitle string, title string, fpath string)
	BarChartWithNameAndNumber       func(dataX []string, dataY []float64, ytitle string, title string, fpath string)
	PieChartWithNameAndNumber       func(dataX []string, dataY []float64, title string, fpath string)
}

var chartstruct chartStruct

func init() {
	chartstruct = chartStruct{
		LineChartWithTimestampAndNumber: drawLineChartWithTimeSeries,
		LineChartWithNumberAndNumber:    drawLineChartWithNumberSeries,
		BarChartWithNameAndNumber:       drawBarChartWithNumberSeries,
		PieChartWithNameAndNumber:       drawPieChartWithNumberSeries,
	}
}

func drawLineChartWithTimeSeries(timestampX []int64, dataY []float64, xtitle string, ytitle string, title string, fpath string) {
	if len(timestampX) != len(dataY) {
		panic(newerr("The number of elements on the X axis and Y axis must be the same"))
	}

	var dataX []time.Time
	for _, i := range timestampX {
		tm := time.Unix(i, 0)
		dataX = append(dataX, tm)
	}

	// statikFS, err := fs.New()
	// panicerr(err)

	// simHeiFile, err := statikFS.Open("/SimHei.ttf")
	// panicerr(err)
	// fontBytes, err := ioutil.ReadAll(simHeiFile)
	// panicerr(err)
	// font, err := truetype.Parse(fontBytes)
	// panicerr(err)
	// simHeiFile.Close()

	graph := chart.Chart{
		//Font:  font,
		Title: title,
		TitleStyle: chart.Style{
			Show: true,
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top: 60,
			},
		},
		Height: 768,
		Width:  2000,
		XAxis: chart.XAxis{
			Name: xtitle,
			NameStyle: chart.Style{
				Show: true,
			},
			Style: chart.Style{
				Show: true,
			},
		},
		YAxis: chart.YAxis{
			Name: ytitle,
			NameStyle: chart.Style{
				Show: true,
			},
			Style: chart.Style{
				Show: true,
			},
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Style: chart.Style{
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
					Show:        true,
				},
				XValues: dataX,
				YValues: dataY,
			},
		},
	}

	f, err := os.Create(fpath)
	panicerr(err)
	defer f.Close()
	graph.Render(chart.PNG, f)
}

func drawLineChartWithNumberSeries(dataX []float64, dataY []float64, xtitle string, ytitle string, title string, fpath string) {
	if len(dataX) != len(dataY) {
		panic(newerr("The number of elements on the X axis and Y axis must be the same"))
	}

	// statikFS, err := fs.New()
	// panicerr(err)

	// simHeiFile, err := statikFS.Open("/SimHei.ttf")
	// panicerr(err)
	// fontBytes, err := ioutil.ReadAll(simHeiFile)
	// panicerr(err)
	// font, err := truetype.Parse(fontBytes)
	// panicerr(err)
	// simHeiFile.Close()

	graph := chart.Chart{
		// Font:  font,
		Title: title,
		TitleStyle: chart.Style{
			Show: true,
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top: 60,
			},
		},
		Height: 768,
		Width:  2000,
		XAxis: chart.XAxis{
			Name: xtitle,
			Style: chart.Style{
				Show: true,
			},
			NameStyle: chart.Style{
				Show: true,
			},
		},
		YAxis: chart.YAxis{
			Name: ytitle,
			NameStyle: chart.Style{
				Show: true,
			},
			Style: chart.Style{
				Show: true,
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
				},
				XValues: dataX,
				YValues: dataY,
			},
		},
	}

	f, err := os.Create(fpath)
	panicerr(err)
	defer f.Close()
	graph.Render(chart.PNG, f)
}

func drawBarChartWithNumberSeries(dataX []string, dataY []float64, ytitle string, title string, fpath string) {
	if len(dataX) != len(dataY) {
		panic(newerr("The number of elements on the X axis and Y axis must be the same"))
	}

	// statikFS, err := fs.New()
	// panicerr(err)

	// simHeiFile, err := statikFS.Open("/SimHei.ttf")
	// panicerr(err)
	// fontBytes, err := ioutil.ReadAll(simHeiFile)
	// panicerr(err)
	// font, err := truetype.Parse(fontBytes)
	// panicerr(err)
	// simHeiFile.Close()

	var chartValue []chart.Value
	for i := 0; i < len(dataX); i++ {
		chartValue = append(chartValue, chart.Value{Value: dataY[i], Label: dataX[i]})
	}

	graph := chart.BarChart{
		// Font:  font,
		Title: title,
		TitleStyle: chart.Style{
			Show: true,
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top: 60,
			},
		},
		YAxis: chart.YAxis{
			Name: ytitle,
			NameStyle: chart.Style{
				Show: true,
			},
			Style: chart.Style{
				Show: true,
			},
		},
		Height: 768,
		Width:  len(dataX) * 120,
		Bars:   chartValue,
	}

	f, err := os.Create(fpath)
	panicerr(err)
	defer f.Close()
	graph.Render(chart.PNG, f)
}

func drawPieChartWithNumberSeries(dataX []string, dataY []float64, title string, fpath string) {
	if len(dataX) != len(dataY) {
		panic(newerr("The number of elements on the X axis and Y axis must be the same"))
	}

	// statikFS, err := fs.New()
	// panicerr(err)

	// simHeiFile, err := statikFS.Open("/SimHei.ttf")
	// panicerr(err)
	// fontBytes, err := ioutil.ReadAll(simHeiFile)
	// panicerr(err)
	// font, err := truetype.Parse(fontBytes)
	// panicerr(err)
	// simHeiFile.Close()

	var chartValue []chart.Value
	for i := 0; i < len(dataX); i++ {
		chartValue = append(chartValue, chart.Value{Value: dataY[i], Label: dataX[i]})
	}

	graph := chart.PieChart{
		// Font:   font,
		Height: 2000,
		Width:  2000,
		Values: chartValue,
		Title:  title,
		TitleStyle: chart.Style{
			Show: true,
		},
	}

	f, err := os.Create(fpath)
	panicerr(err)
	defer f.Close()
	graph.Render(chart.PNG, f)
}
