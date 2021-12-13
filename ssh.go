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
	Panicerr(err)

	return &sshStruct{
		client: client,
		user:   user,
		pass:   pass,
		port:   port,
		host:   host,
	}
}

func (m *sshStruct) Close() {
	m.client.Close()
}

func (m *sshStruct) Exec(cmd string) (output string, status int) {
	session, err := m.client.NewSession()
	if err != nil {
		m.client.Close()
		Panicerr(err)
	}

	out, err := session.CombinedOutput(cmd)
	output = string(out)
	//lg.debug(output, err)
	if err != nil {
		if String("Process exited with status ").In(err.Error()) {
			o := String(err.Error()).Split()
			status = Int(o[len(o)-1])
		} else {
			Panicerr(err)
		}
	} else {
		status = 0
	}
	return
}

func (m *sshStruct) PushFile(local string, remote string) {
	sftp, err := sftp.NewClient(m.client)
	Panicerr(err)
	defer sftp.Close()

	fr, err := sftp.Create(remote)
	Panicerr(err)

	fl, err := os.Open(local)
	Panicerr(err)

	io.Copy(fr, fl)

	err = fr.Close()
	Panicerr(err)

	err = fl.Close()
	Panicerr(err)
}

func (m *sshStruct) PullFile(remote string, local string) {
	sftp, err := sftp.NewClient(m.client)
	Panicerr(err)
	defer sftp.Close()

	fr, err := sftp.Open(remote)
	Panicerr(err)

	fl, err := os.Create(local)
	Panicerr(err)

	io.Copy(fl, fr)

	err = fr.Close()
	Panicerr(err)

	err = fl.Close()
	Panicerr(err)
}
