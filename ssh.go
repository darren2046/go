package golanglibs

import (
	"io"
	"net"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type sshStruct struct {
	client  *ssh.Client
	session *ssh.Session
	user    string
	pass    string
	host    string
	port    int
}

func getSSH(user string, pass string, host string, port int) *sshStruct {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(pass)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", host+":"+Str(port), sshConfig)
	panicerr(err)

	return &sshStruct{
		client: client,
		user:   user,
		pass:   pass,
		port:   port,
		host:   host,
	}
}

func (m *sshStruct) close() {
	m.client.Close()
}

func (m *sshStruct) exec(cmd string) (output string, status int) {
	session, err := m.client.NewSession()
	if err != nil {
		m.client.Close()
		panicerr(err)
	}

	out, err := session.CombinedOutput(cmd)
	output = string(out)
	//lg.debug(output, err)
	if err != nil {
		if String("Process exited with status ").In(err.Error()) {
			o := String(err.Error()).Split()
			status = Int(o[len(o)-1])
		} else {
			panicerr(err)
		}
	} else {
		status = 0
	}
	return
}

func (m *sshStruct) pushFile(local string, remote string) {
	sftp, err := sftp.NewClient(m.client)
	panicerr(err)
	defer sftp.Close()

	fr, err := sftp.Create(remote)
	panicerr(err)

	fl, err := os.Open(local)
	panicerr(err)

	io.Copy(fr, fl)

	err = fr.Close()
	panicerr(err)

	err = fl.Close()
	panicerr(err)
}

func (m *sshStruct) pullFile(remote string, local string) {
	sftp, err := sftp.NewClient(m.client)
	panicerr(err)
	defer sftp.Close()

	fr, err := sftp.Open(remote)
	panicerr(err)

	fl, err := os.Create(local)
	panicerr(err)

	io.Copy(fl, fr)

	err = fr.Close()
	panicerr(err)

	err = fl.Close()
	panicerr(err)
}
