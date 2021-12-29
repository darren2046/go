package golanglibs

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

type cmdStruct struct {
	GetOutput                func(command string, timeoutSecond ...interface{}) string
	GetStatusOutput          func(command string, timeoutSecond ...interface{}) (int, string)
	GetOutputWithShell       func(command string, timeoutSecond ...interface{}) string
	GetStatusOutputWithShell func(command string, timeoutSecond ...interface{}) (int, string)
	Tail                     func(command string) chan string
	Exists                   func(cmd string) bool
	Which                    func(cmd string) (path string)
}

var Cmd cmdStruct

func init() {
	Cmd = cmdStruct{
		GetOutput:                getOutput,
		GetStatusOutput:          getStatusOutput,
		GetOutputWithShell:       getOutputWithShell,
		GetStatusOutputWithShell: getStatusOutputWithShell,
		Tail:                     tailCmdOutput,
		Exists:                   cmdExists,
		Which:                    cmdWhich,
	}
}

func cmdWhich(cmd string) (path string) {
	path, err := exec.LookPath(cmd)
	if err != nil {
		Panicerr("Command not found: " + cmd)
	}
	return
}

func cmdExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

type execCmdStruct struct {
	buf string
	out chan string
}

func (m *execCmdStruct) Write(p []byte) (n int, err error) {
	for _, pp := range p {
		ps := string(pp)
		if ps == "\n" {
			Try(func() {
				m.out <- m.buf
			})
			m.buf = ""
		} else {
			m.buf = m.buf + ps
		}
	}

	return len(p), nil
}

func tailCmdOutput(command string) chan string {
	w := execCmdStruct{out: make(chan string, 9999)}
	go func() {
		cmd := exec.Command("/bin/bash", "-c", command)
		cmd.Stdin = os.Stdin
		cmd.Stdout = &w
		cmd.Stderr = &w
		err := cmd.Start()
		if err != nil {
			close(w.out)
			Panicerr(err)
		}
		err = cmd.Wait()
		close(w.out)
		if err != nil {
			Panicerr(err)
		}
	}()
	return w.out
}

func getOutputWithShell(command string, timeoutSecond ...interface{}) string {
	_, o := getStatusOutputWithShell(command, timeoutSecond...)
	return o
}

// subprocess.getstautsoutput()
// command.getstatusoutput()
func getStatusOutputWithShell(command string, timeoutSecond ...interface{}) (int, string) {
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

	if !Cmd.Exists(parts[0]) {
		Panicerr("Command not exists")
	}

	var shell string
	for _, s := range []string{"/bin/bash"} {
		if Os.Path.Exists(s) {
			shell = s
			break
		}
	}

	if shell == "" {
		Panicerr("Shell not found")
	}

	var statuscode int
	var output string
	if len(timeoutSecond) != 0 {
		t := timeoutSecond[0]
		timeoutDuration := getTimeDuration(t)
		ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
		defer cancel()

		cmd := exec.CommandContext(ctx, shell, "-c", command)
		//cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)

		out, err := cmd.CombinedOutput()

		if err != nil {
			e := err.(*exec.ExitError)
			statuscode = e.ExitCode()
		} else {
			statuscode = 0
		}
		output = string(out)
	} else {
		cmd := exec.Command(shell, "-c", command)

		out, err := cmd.CombinedOutput()

		if err != nil {
			e := err.(*exec.ExitError)
			statuscode = e.ExitCode()
		} else {
			statuscode = 0
		}
		output = string(out)
	}
	return statuscode, output
}

func getOutput(command string, timeoutSecond ...interface{}) string {
	_, o := getStatusOutput(command, timeoutSecond...)
	return o
}

// subprocess.getstautsoutput()
// command.getstatusoutput()
func getStatusOutput(command string, timeoutSecond ...interface{}) (int, string) {
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

	if !Cmd.Exists(parts[0]) {
		Panicerr("Command not exists")
	}

	var statuscode int
	var output string
	if len(timeoutSecond) != 0 {
		t := timeoutSecond[0]
		timeoutDuration := getTimeDuration(t)
		ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
		defer cancel()

		// cmd := exec.CommandContext(ctx, "/bin/bash", "-c", command)
		cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)

		out, err := cmd.CombinedOutput()

		if err != nil {
			e := err.(*exec.ExitError)
			statuscode = e.ExitCode()
		} else {
			statuscode = 0
		}
		output = string(out)
	} else {
		// cmd := exec.Command("/bin/sh", "-c", command)
		cmd := exec.Command(parts[0], parts[1:]...)

		out, err := cmd.CombinedOutput()

		if err != nil {
			e := err.(*exec.ExitError)
			statuscode = e.ExitCode()
		} else {
			statuscode = 0
		}
		output = string(out)
	}
	return statuscode, output
}
