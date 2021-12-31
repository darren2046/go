package golanglibs

import (
	"os"
	"os/exec"
	"strings"
	"unicode"

	"github.com/creack/pty"
)

type pexpectStruct struct {
	cmd         *exec.Cmd
	bufall      string // 所有屏幕的显示内容，包括了输入
	ptmx        *os.File
	logToStdout bool // 是否在屏幕打印出整个交互（适合做debug)
	isAlive     bool // 子进程是否在运行
}

func (m *pexpectStruct) Sendline(msg string) {
	_, err := m.ptmx.Write([]byte(msg + "\n"))
	Panicerr(err)
}

func (m *pexpectStruct) Close() {
	m.isAlive = false
	m.ptmx.Close()
	m.cmd.Process.Signal(os.Kill)
	m.cmd.Wait()
}

func (m *pexpectStruct) ExitCode() int {
	return m.cmd.ProcessState.ExitCode()
}

func (m *pexpectStruct) IsAlive() bool {
	return m.isAlive
}

func (m *pexpectStruct) LogToStdout(enable ...bool) {
	var e bool
	if len(enable) != 0 {
		e = enable[0]
	} else {
		e = true
	}
	m.logToStdout = e
}

func (m *pexpectStruct) GetLog() *stringStruct {
	return String(m.bufall)
}

func (m *pexpectStruct) ClearLog() {
	m.bufall = ""
}

func (m *pexpectStruct) Wait() {
	m.cmd.Wait()
}

func pexpect(command string) *pexpectStruct {
	q := rune(0)
	parts := strings.FieldsFunc(command, func(r rune) bool {
		switch {
		case r == q:
			q = rune(0)
			return false
		case q != rune(0):
			return false
		case unicode.In(r, unicode.Quotation_Mark):
			q = r
			return false
		default:
			return unicode.IsSpace(r)
		}
	})
	// remove the " and ' on both sides
	for i, v := range parts {
		f, l := v[0], len(v)
		if l >= 2 && (f == '"' || f == '\'') {
			parts[i] = v[1 : l-1]
		}
	}

	m := pexpectStruct{
		isAlive: true,
	}

	if !cmdExists(parts[0]) {
		Panicerr("Command not exists")
	}

	cmd := exec.Command(parts[0], parts[1:]...)

	m.cmd = cmd

	var err error
	m.ptmx, err = pty.StartWithSize(cmd, &pty.Winsize{
		Rows: 60,
		Cols: 1024,
	})
	Panicerr(err)

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := m.ptmx.Read(buf)
			if err != nil {
				break
			}

			m.bufall = m.bufall + string(buf[:n])
			if m.logToStdout {
				os.Stdout.Write([]byte(buf[:n]))
			}
		}

		m.Close()
		m.isAlive = false
	}()

	return &m
}
