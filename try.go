package golanglibs

import (
	"errors"
	"fmt"
	"runtime/debug"
)

type exception struct {
	Error error
}

type tryConfig struct {
	retry int // retry times while error occure
	sleep int // sleep seconds between retry
}

func throw() {
	panic("_____rethrow")
}

func try(f func(), trycfg ...tryConfig) (e exception) {
	if len(trycfg) == 0 {
		e = exception{nil}
		defer func() {
			if err := recover(); err != nil {
				errmsg := fmt.Sprintf("%s", err)
				if len(Re.FindAll(".+\\.go:[0-9]+ >> .+? >> \\(.+?\\)", errmsg)) == 0 {
					e.Error = errors.New(fmtDebugStack(errmsg, str(debug.Stack())))
				} else {
					e.Error = errors.New(errmsg)
				}
			}
		}()
		f()
		return
	}
	for i := 0; ; i++ {
		e = func() (e exception) {
			e = exception{nil}
			defer func() {
				if err := recover(); err != nil {
					e.Error = errors.New(fmt.Sprintf("%s", err))
				}
			}()
			f()
			return
		}()
		if e.Error == nil {
			return
		}
		if e.Error != nil && trycfg[0].retry > 0 && i >= trycfg[0].retry {
			break
		}
		Time.Sleep(trycfg[0].sleep)
	}
	return
}

func (e exception) except(f func(err error)) error {
	if e.Error != nil {
		defer func() {
			if err := recover(); err != nil {
				if err == "_____rethrow" {
					err = e.Error
				}
				panic(err)
			}
		}()
		f(e.Error)
	}
	return e.Error
}
