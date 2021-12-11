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
	buf         string // sendline之前的命令输出，每次sendline都会清空
	bufall      string // 所有屏幕的显示内容，包括了输入
	ptmx        *os.File
	logToStdout bool // 是否在屏幕打印出整个交互（适合做debug)
	isAlive     bool // 子进程是否在运行
}

func (m *pexpectStruct) sendline(msg string) {
	m.buf = ""
	m.bufall += msg + "\n"
	_, err := m.ptmx.Write([]byte(msg + "\n"))
	panicerr(err)
}

func (m *pexpectStruct) close() {
	m.isAlive = false
	m.ptmx.Close()
	m.cmd.Process.Signal(os.Kill)
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
		panicerr("Command not exists")
	}

	cmd := exec.Command(parts[0], parts[1:]...)

	m.cmd = cmd

	var err error
	m.ptmx, err = pty.StartWithSize(cmd, &pty.Winsize{
		Rows: 60,
		Cols: 1024,
	})
	panicerr(err)

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := m.ptmx.Read(buf)
			if err != nil {
				break
			}

			m.buf = m.buf + string(buf[:n])
			m.bufall += string(buf[:n])
			if m.logToStdout {
				os.Stdout.Write([]byte(buf[:n]))
			}
		}

		m.close()
		m.isAlive = false
	}()

	return &m
}
