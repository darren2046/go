package golanglibs

import "testing"

func TestUnixSocket(t *testing.T) {
	Print("Start")
	sp := "/tmp/kkkwjjjskkkwjefwehfliwhegliwehg.socks"
	go func() {
		ss := Socket.UNIX.Listen(sp)
		defer ss.Close()
		cs := <-ss.Accept()
		Lg.Trace("Send to client: c1")
		cs.Send("c1")
		Lg.Trace("Recv from client:", cs.Recv(1024))
	}()
	Time.Sleep(1)
	cc := Socket.UNIX.Connect(sp)
	Lg.Trace("Recv from server:", cc.Recv(1024))
	Lg.Trace("Send to server: c2")
	cc.Send("c2")
	cc.Close()
	Time.Sleep(1)
}
