package golanglibs

import (
	"math"
	"strconv"
	"time"
)

type timeStruct struct {
	Now            func() float64
	TimeDuration   func(seconds interface{}) time.Duration
	FormatDuration func(second int64) (result string)
	Sleep          func(t interface{})
}

var Time timeStruct

func init() {
	Time = timeStruct{
		Now:            timeNowInTimestamp,
		TimeDuration:   getTimeDuration,
		FormatDuration: fmtTimeDuration,
		Sleep:          sleep,
	}
}

func sleep(t interface{}) {
	time.Sleep(Time.TimeDuration(t))
}

func plural(count int, singular string) (result string) {
	if (count == 1) || (count == 0) {
		result = strconv.Itoa(count) + " " + singular + " "
	} else {
		result = strconv.Itoa(count) + " " + singular + "s "
	}
	return
}

func fmtTimeDuration(second int64) (result string) {
	years := math.Floor(float64(second) / 60 / 60 / 24 / 7 / 30 / 12)
	seconds := second % (60 * 60 * 24 * 7 * 30 * 12)
	months := math.Floor(float64(seconds) / 60 / 60 / 24 / 7 / 30)
	seconds = second % (60 * 60 * 24 * 7 * 30)
	weeks := math.Floor(float64(seconds) / 60 / 60 / 24 / 7)
	seconds = second % (60 * 60 * 24 * 7)
	days := math.Floor(float64(seconds) / 60 / 60 / 24)
	seconds = second % (60 * 60 * 24)
	hours := math.Floor(float64(seconds) / 60 / 60)
	seconds = second % (60 * 60)
	minutes := math.Floor(float64(seconds) / 60)
	seconds = second % 60

	if years > 0 {
		result = plural(int(years), "year") + plural(int(months), "month") + plural(int(weeks), "week") + plural(int(days), "day") + plural(int(hours), "hour") + plural(int(minutes), "minute") + plural(int(seconds), "second")
	} else if months > 0 {
		result = plural(int(months), "month") + plural(int(weeks), "week") + plural(int(days), "day") + plural(int(hours), "hour") + plural(int(minutes), "minute") + plural(int(seconds), "second")
	} else if weeks > 0 {
		result = plural(int(weeks), "week") + plural(int(days), "day") + plural(int(hours), "hour") + plural(int(minutes), "minute") + plural(int(seconds), "second")
	} else if days > 0 {
		result = plural(int(days), "day") + plural(int(hours), "hour") + plural(int(minutes), "minute") + plural(int(seconds), "second")
	} else if hours > 0 {
		result = plural(int(hours), "hour") + plural(int(minutes), "minute") + plural(int(seconds), "second")
	} else if minutes > 0 {
		result = plural(int(minutes), "minute") + plural(int(seconds), "second")
	} else {
		result = plural(int(seconds), "second")
	}

	result = String(result).Strip().Get()

	return
}

func getTimeDuration(seconds interface{}) time.Duration {
	var timeDuration time.Duration
	if typeof(seconds) == "float64" {
		tt := seconds.(float64) * 1000
		if tt < 0 {
			tt = 0
		}
		timeDuration = time.Duration(tt) * time.Millisecond
	}
	if typeof(seconds) == "int" || typeof(seconds) == "int8" || typeof(seconds) == "int16" || typeof(seconds) == "int32" || typeof(seconds) == "int64" {
		tt := toInt64(seconds)
		if tt < 0 {
			tt = 0
		}
		timeDuration = time.Duration(tt) * time.Second
	}
	return timeDuration
}

func timeNowInTimestamp() float64 {
	return toFloat64(time.Now().UnixMicro()) / 1000000
}
