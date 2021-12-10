package golanglibs

import (
	"fmt"
	"runtime/debug"
	"strings"
)

type errorStruct struct {
	msg string
}

func fmtDebugStack(msg string, stack string) string {
	//lg.debug("msg:", msg)
	//lg.debug("stack:", stack)

	blackFileList := []string{
		"lib.go",
		"stack.go",
	}

	l := Re.FindAll("([\\-a-zA-Z0-9]+\\.go:[0-9]+)", stack)
	//lg.debug(l)
	for i, j := 0, len(l)-1; i < j; i, j = i+1, j-1 {
		l[i], l[j] = l[j], l[i]
	}
	//lg.debug(l)

	var link []string
	for _, f := range l {
		ff := strings.Split(f[0], ":")[0]
		inside := func(a string, list []string) bool {
			for _, b := range list {
				if b == a {
					return true
				}
			}
			return false
		}(ff, blackFileList)
		if !inside {
			link = append(link, f[0])
		}
	}
	//lg.debug(link)

	var strr string
	if len(link) != 1 {
		// strr = link[len(link)-2] + " >> " + msg + " >> " + "(" + strJoin(" => ", link[:len(link)-1]) + ")"
		strr = link[len(link)-1] + " >> " + msg + " >> " + "(" + String(" => ").Join(link).Get() + ")"
	} else {
		strr = link[0] + " >> " + msg
	}

	//lg.debug("strr:", strr)
	return strr
}

func newerr(msg interface{}) *errorStruct {
	switch t := msg.(type) {
	case string:
		return &errorStruct{
			msg: fmtDebugStack(t, string(debug.Stack())),
		}
	case error:
		return &errorStruct{
			msg: fmtDebugStack(t.Error(), string(debug.Stack())),
		}
	default:
		return &errorStruct{
			msg: fmtDebugStack(fmt.Sprintf("%s", t), string(debug.Stack())),
		}
	}
}

func panicerr(err interface{}) {
	switch t := err.(type) {
	case string:
		// lg.trace("1")
		panic(fmtDebugStack(t, string(debug.Stack())))
	case error:
		// lg.trace("2")
		panic(fmtDebugStack(t.Error(), string(debug.Stack())))
	case *errorStruct:
		// lg.trace(3)
		panic(t.msg)
	case nil:
		return
	default:
		panic(fmtDebugStack(fmt.Sprintf("%s", t), string(debug.Stack())))
	}
}
