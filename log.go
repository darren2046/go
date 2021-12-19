package golanglibs

import (
	"fmt"
	"runtime"

	"github.com/fatih/color"
)

var Lg *logStruct

func init() {
	Lg = getLogger()
}

type logStruct struct {
	level                    []string
	levelString              string
	json                     bool
	color                    bool
	logDir                   string
	logFileName              string
	logFileSuffix            string
	fd                       *fileIOStruct
	displayOnTerminal        bool
	lock                     *lockStruct
	logfiles                 []string
	maxlogfiles              int
	logFileSizeInMB          int
	currentLogFileSizeInByte int
	currentLogFileNumber     int
}

func getLogger() *logStruct {
	return &logStruct{
		level:                    []string{"TRAC", "DEBU", "INFO", "WARN", "ERRO"},
		color:                    true,
		displayOnTerminal:        true,
		lock:                     getLock(),
		logFileSizeInMB:          0,
		currentLogFileSizeInByte: 0,
		currentLogFileNumber:     0,
	}
}

func (m *logStruct) SetLevel(level string) {
	if level == "trace" {
		m.level = []string{"TRAC", "DEBU", "INFO", "WARN", "ERRO"}
	} else if level == "debug" {
		m.level = []string{"DEBU", "INFO", "WARN", "ERRO"}
	} else if level == "info" {
		m.level = []string{"INFO", "WARN", "ERRO"}
	} else if level == "warn" {
		m.level = []string{"WARN", "ERRO"}
	} else if level == "error" {
		m.level = []string{"ERRO"}
	} else if level == "" {
		m.level = []string{}
	}
	m.levelString = level
}

func (m *logStruct) GetLevel() string {
	return m.levelString
}

func (m *logStruct) SetLogFile(path string, maxLogFileCount int, logFileSizeInMB ...int) {
	m.logDir = Os.Path.Basedir(path)
	if !pathExists(m.logDir) {
		mkdir(m.logDir)
	}

	f := String(Os.Path.Basename(path)).Split(".")
	m.logFileName = String(".").Join(f[:len(f)-1]).Get()
	m.logFileSuffix = f[len(f)-1]

	var logpath string
	if len(logFileSizeInMB) != 0 {
		m.logFileSizeInMB = logFileSizeInMB[0]
		for {
			logpath = pathJoin(m.logDir, m.logFileName+"."+Str(m.currentLogFileNumber)+"."+m.logFileSuffix)
			if pathExists(logpath) {
				m.currentLogFileNumber++
			} else {
				break
			}
		}
	} else {
		logpath = pathJoin(m.logDir, m.logFileName+"."+strftime("%Y-%m-%d", Time.Now())+"."+m.logFileSuffix)
	}
	m.fd = Open(logpath, "a")
	m.logfiles = append(m.logfiles, logpath)

	m.maxlogfiles = maxLogFileCount
}

func (m *logStruct) Error(args ...interface{}) {
	t := strftime("%m-%d %H:%M:%S", Time.Now())
	level := "ERRO"

	var msgarr []string
	for _, a := range args {
		msgarr = append(msgarr, fmt.Sprint(a))
	}
	msg := String(" ").Join(msgarr).Get()

	_, file, no, _ := runtime.Caller(1)
	position := Os.Path.Basename(file) + ":" + Str(no)

	m.show(t, level, msg, position)
}

func (m *logStruct) Warn(args ...interface{}) {
	t := strftime("%m-%d %H:%M:%S", Time.Now())
	level := "WARN"

	var msgarr []string
	for _, a := range args {
		msgarr = append(msgarr, fmt.Sprint(a))
	}
	msg := String(" ").Join(msgarr).Get()

	_, file, no, _ := runtime.Caller(1)
	position := Os.Path.Basename(file) + ":" + Str(no)

	m.show(t, level, msg, position)
}

func (m *logStruct) Info(args ...interface{}) {
	t := strftime("%m-%d %H:%M:%S", Time.Now())
	level := "INFO"

	var msgarr []string
	for _, a := range args {
		msgarr = append(msgarr, fmt.Sprint(a))
	}
	msg := String(" ").Join(msgarr).Get()

	_, file, no, _ := runtime.Caller(1)
	Print(file, no)
	position := Os.Path.Basename(file) + ":" + Str(no)

	m.show(t, level, msg, position)
}

func (m *logStruct) Trace(args ...interface{}) {
	t := strftime("%m-%d %H:%M:%S", Time.Now())
	level := "TRAC"

	var msgarr []string
	for _, a := range args {
		msgarr = append(msgarr, fmt.Sprint(a))
	}
	msg := String(" ").Join(msgarr).Get()

	_, file, no, _ := runtime.Caller(1)
	position := Os.Path.Basename(file) + ":" + Str(no)

	m.show(t, level, msg, position)
}

func (m *logStruct) Debug(args ...interface{}) {
	t := strftime("%m-%d %H:%M:%S", Time.Now())
	level := "DEBU"

	var msgarr []string
	for _, a := range args {
		msgarr = append(msgarr, Sprint(a))
	}
	msg := String(" ").Join(msgarr).Get()

	_, file, no, _ := runtime.Caller(1)
	position := Os.Path.Basename(file) + ":" + Str(no)

	m.show(t, level, msg, position)
}

func (m *logStruct) show(t string, level string, msg string, position string) {
	if Array(m.level).Has(level) {
		var strMsg string
		if m.json {
			strMsg = jsonDumps(map[string]string{
				"time":    t,
				"level":   level,
				"message": msg,
			})
		} else {
			msg = String(msg).Replace("\n", "\n                      ").Get()
			if m.color {
				if level == "ERRO" {
					strMsg = color.RedString(t + fmt.Sprintf(" %3v", Os.GoroutineID()) + " [" + level + "] (" + position + ") " + msg)
				} else if level == "WARN" {
					strMsg = color.YellowString(t + fmt.Sprintf(" %3v", Os.GoroutineID()) + " [" + level + "] (" + position + ") " + msg)
				} else if level == "INFO" {
					strMsg = color.HiBlueString(t + fmt.Sprintf(" %3v", Os.GoroutineID()) + " [" + level + "] (" + position + ") " + msg)
				} else if level == "TRAC" {
					strMsg = color.MagentaString(t + fmt.Sprintf(" %3v", Os.GoroutineID()) + " [" + level + "] (" + position + ") " + msg)
				} else if level == "DEBU" {
					strMsg = color.HiCyanString(t + fmt.Sprintf(" %3v", Os.GoroutineID()) + " [" + level + "] (" + position + ") " + msg)
				}
			} else {
				strMsg = t + "[" + level + "] (" + position + ") " + msg
			}
		}

		m.lock.Acquire()
		if m.displayOnTerminal {
			fmt.Println(strMsg)
		}
		if m.fd != nil {
			if m.logFileSizeInMB == 0 {
				if m.fd.path != pathJoin(m.logDir, m.logFileName+"."+strftime("%Y-%m-%d", Time.Now())+"."+m.logFileSuffix) {
					m.fd.Close()
					logpath := pathJoin(m.logDir, m.logFileName+"."+strftime("%Y-%m-%d", Time.Now())+"."+m.logFileSuffix)
					m.fd = Open(logpath, "a")
					m.logfiles = append(m.logfiles, logpath)
					if len(m.logfiles) > m.maxlogfiles {
						Os.Unlink(m.logfiles[0])
						m.logfiles = m.logfiles[1:]
					}
				}
			} else {
				if m.currentLogFileSizeInByte > m.logFileSizeInMB*1024*1024 {
					m.currentLogFileSizeInByte = 0
					m.fd.Close()
					var logpath string
					for {
						logpath = pathJoin(m.logDir, m.logFileName+"."+Str(m.currentLogFileNumber)+"."+m.logFileSuffix)
						if pathExists(logpath) {
							m.currentLogFileNumber++
						} else {
							break
						}
					}
					m.fd = Open(logpath, "a")
					m.logfiles = append(m.logfiles, logpath)
					if len(m.logfiles) > m.maxlogfiles {
						Os.Unlink(m.logfiles[0])
						m.logfiles = m.logfiles[1:]
					}
				}
			}
			m.fd.Write(strMsg + "\n")
			m.currentLogFileSizeInByte = m.currentLogFileSizeInByte + len(strMsg) + 1
		}
		m.lock.Release()
	}
}
