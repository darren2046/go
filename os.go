package golanglibs

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"unicode"
)

type osStruct struct {
	Chdir           func(path string)
	Mkdir           func(filename string)
	Getcwd          func() string
	Exit            func(status int)
	Touch           func(filePath string)
	Chmod           func(filePath string, mode uint32) bool
	Chown           func(filePath string, uid, gid int) bool
	Copy            func(filePath, dest string)
	Rename          func(filePath, newPosition string)
	Move            func(filePath, newPosition string)
	Path            *pathStruct
	System          func(command string, timeoutSecond ...interface{}) int
	SystemWithShell func(command string, timeoutSecond ...interface{}) int
	Hostname        func() string
	Envexists       func(varname string) bool
	Getenv          func(varname string) string
	Walk            func(path string) chan string
	Listdir         func(path string) (res []string)
	SelfDir         func() string
	TempFilePath    func() string
	TempDirPath     func() string
	Getuid          func() int
	ProgressAlive   func(pid int) bool
	GoroutineID     func() int64
	Unlink          func(filename string)
}

var Os osStruct

func init() {
	Os = osStruct{
		Chdir:           chdir,
		Mkdir:           mkdir,
		Getcwd:          getcwd,
		Exit:            os.Exit,
		Touch:           touch,
		Chmod:           chmod,
		Chown:           chown,
		Copy:            copyFile,
		Rename:          rename,
		Move:            move,
		Path:            &pathstruct,
		System:          system,
		SystemWithShell: systemWithShell,
		Hostname:        gethostname,
		Envexists:       envexists,
		Getenv:          getenv,
		Walk:            walk,
		Listdir:         listdir,
		SelfDir:         getSelfDir,
		TempFilePath:    getTempFilePath,
		TempDirPath:     getTempDirPath,
		Getuid:          os.Getuid,
		ProgressAlive:   progressAlive,
		GoroutineID:     goroutineID,
		Unlink:          unlink,
	}
}

func unlink(filename string) {
	err := os.RemoveAll(filename)
	Panicerr(err)
}

func goroutineID() int64 {
	var (
		buf [64]byte
		n   = runtime.Stack(buf[:], false)
		stk = strings.TrimPrefix(string(buf[:n]), "goroutine ")
	)

	idField := strings.Fields(stk)[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Errorf("can not get goroutine id: %v", err))
	}

	return int64(id)
}

func progressAlive(pid int) bool {
	if pid <= 0 {
		Panicerr(fmt.Errorf("invalid pid %v", pid))
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = proc.Signal(syscall.Signal(0))
	if err == nil {
		return true
	}
	if err.Error() == "os: process already finished" {
		return false
	}
	errno, ok := err.(syscall.Errno)
	if !ok {
		return false
	}
	switch errno {
	case syscall.ESRCH:
		return false
	case syscall.EPERM:
		return true
	}
	return false
}

// func setFileTime(fpath string, mtime string, atime ...string) {
// 	var at time.Time
// 	if len(atime) == 0 {
// 		fi, err := os.Stat(fpath)
// 		Panicerr(err)
// 		//mtime := fi.ModTime()
// 		stat := fi.Sys().(*syscall.Stat_t)
// 		at = time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec))
// 		// ctime := time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))
// 	} else {
// 		at = time.Unix(strptime("%Y-%m-%d %H:%M:%S", atime[0]), 0)
// 	}

// 	err := os.Chtimes(fpath, at, time.Unix(strptime("%Y-%m-%d %H:%M:%S", mtime), 0))
// 	Panicerr(err)
// }

func getTempDirPath() string {
	dir, err := ioutil.TempDir("", "systemd-private-")
	Panicerr(err)
	return dir
}

func getTempFilePath() string {
	file, err := ioutil.TempFile("", "systemd-private-")
	Panicerr(err)
	defer os.Remove(file.Name())

	return file.Name()
}

func getSelfDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}

func listdir(path string) (res []string) {
	files, err := ioutil.ReadDir(path)
	Panicerr(err)

	for _, f := range files {
		res = append(res, f.Name())
	}
	return
}

func walk(path string) chan string {
	pc := make(chan string)
	go func() {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if path != "." && path != ".." {
				pc <- path
			}
			return nil
		})
		if err != nil {
			close(pc)
			_, fn, line, _ := runtime.Caller(0)
			panic(filepath.Base(fn) + ":" + strconv.Itoa(line-9) + " >> " + err.Error())
		}
		close(pc)
	}()

	return pc
}

func getenv(varname string) string {
	e, exist := os.LookupEnv(varname)
	if !exist {
		Panicerr("Env not exists")
	}
	return e
}

func envexists(varname string) bool {
	_, exist := os.LookupEnv(varname)
	return exist
}

func gethostname() string {
	name, err := os.Hostname()
	Panicerr(err)
	return name
}

func systemWithShell(command string, timeoutSecond ...interface{}) int {
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
	if len(timeoutSecond) != 0 {
		t := timeoutSecond[0]
		timeoutDuration := getTimeDuration(t)
		ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
		defer cancel()

		cmd := exec.CommandContext(ctx, shell, "-c", command)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		if err != nil {
			e := err.(*exec.ExitError)
			statuscode = e.ExitCode()
		} else {
			statuscode = 0
		}
	} else {
		cmd := exec.Command(shell, "-c", command)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		if err != nil {
			e := err.(*exec.ExitError)
			statuscode = e.ExitCode()
		} else {
			statuscode = 0
		}
	}
	return statuscode
}

func system(command string, timeoutSecond ...interface{}) int {
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
	if len(timeoutSecond) != 0 {
		t := timeoutSecond[0]
		timeoutDuration := getTimeDuration(t)
		ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
		defer cancel()

		// cmd := exec.CommandContext(ctx, "/bin/bash", "-c", command) // 如果不是bash会kill不掉
		cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		cmd.Wait()

		if err != nil {
			e := err.(*exec.ExitError)
			statuscode = e.ExitCode()
		} else {
			statuscode = 0
		}
	} else {
		// cmd := exec.Command("/bin/bash", "-c", command)
		cmd := exec.Command(parts[0], parts[1:]...)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		cmd.Wait()

		if err != nil {
			e := err.(*exec.ExitError)
			statuscode = e.ExitCode()
		} else {
			statuscode = 0
		}
	}
	return statuscode
}

func copyFile(filePath, dest string) {
	fd1, err := os.Open(filePath)
	Panicerr(err)
	defer fd1.Close()
	fd2, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0644)
	Panicerr(err)
	defer fd2.Close()
	_, err = io.Copy(fd2, fd1)
	Panicerr(err)
}

func rename(filePath, newPosition string) {
	err := os.Rename(filePath, newPosition)
	Panicerr(err)
}

func move(filePath, newPosition string) {
	err := os.Rename(filePath, newPosition)
	Panicerr(err)
}

func mkdir(filename string) {
	err := os.MkdirAll(filename, 0755)
	Panicerr(err)
}

func getcwd() string {
	dir, err := os.Getwd()
	Panicerr(err)
	return dir
}

func touch(filePath string) {
	fd, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0666)
	Panicerr(err)
	fd.Close()
}

func chmod(filePath string, mode uint32) bool {
	if len(Str(mode)) == 3 {
		mode = Uint32("0" + Str(mode))
	}
	return os.Chmod(filePath, os.FileMode(mode)) == nil
}

func chown(filePath string, uid, gid int) bool {
	return os.Chown(filePath, uid, gid) == nil
}

func chdir(path string) {
	err := os.Chdir(path)
	Panicerr(err)
}
