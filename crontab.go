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

func (m *crontabStruct) Add(schedule string, fn interface{}, args ...interface{}) {
	err := m.c.AddJob(schedule, fn, args...)
	Panicerr(err)
}

func (m *crontabStruct) Destory() {
	m.c.Shutdown()
}
