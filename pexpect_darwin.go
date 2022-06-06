package golanglibs

import (
	"os"
	"os/exec"
	"strings"
	"unicode"

	"github.com/creack/pty"
)

func (m *PexpectStruct) Sendline(msg string) {
	_, err := m.Ptmx.Write([]byte(msg + "\n"))
	Panicerr(err)
}

func (m *PexpectStruct) Close() {
	m.isAlive = false
	m.Ptmx.Close()
	m.Cmd.Process.Signal(os.Kill)
	m.Cmd.Wait()
}

func (m *PexpectStruct) ExitCode() int {
	return m.Cmd.ProcessState.ExitCode()
}

func (m *PexpectStruct) IsAlive() bool {
	return m.isAlive
}

func (m *PexpectStruct) LogToStdout(enable ...bool) {
	var e bool
	if len(enable) != 0 {
		e = enable[0]
	} else {
		e = true
	}
	m.logToStdout = e
}

func (m *PexpectStruct) GetLog() *stringStruct {
	return String(m.bufall)
}

func (m *PexpectStruct) ClearLog() {
	m.bufall = ""
}

func (m *PexpectStruct) Wait() {
	m.Cmd.Wait()
}

func pexpect(command string) *PexpectStruct {
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

	m := PexpectStruct{
		isAlive: true,
	}

	if !CmdExists(parts[0]) {
		Panicerr("Command not exists")
	}

	Cmd := exec.Command(parts[0], parts[1:]...)

	m.Cmd = Cmd

	var err error
	m.Ptmx, err = pty.StartWithSize(Cmd, &pty.Winsize{
		Rows: 60,
		Cols: 1024,
	})
	Panicerr(err)

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := m.Ptmx.Read(buf)
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
