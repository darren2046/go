package golanglibs

import (
	"github.com/mileusna/crontab"
)

type crontabStruct struct {
	c *crontab.Crontab
}

func getCrontab() *crontabStruct {
	return &crontabStruct{c: crontab.New()}
}

func (m *crontabStruct) add(schedule string, fn interface{}, args ...interface{}) {
	err := m.c.AddJob(schedule, fn, args...)
	panicerr(err)
}

func (m *crontabStruct) destory() {
	m.c.Shutdown()
}
