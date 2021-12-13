package golanglibs

import (
	"encoding/binary"
	"io"
	"math/rand"
	"time"

	"github.com/xtaci/smux"
)

type smuxStruct struct {
	ServerWrapper func(conn io.ReadWriteCloser, cfg ...SmuxConfig) *smuxServerSideListener
	ClientWrapper func(conn io.ReadWriteCloser, cfg ...SmuxConfig) *smuxClientSideSession
}

var smuxstruct smuxStruct

func init() {
	smuxstruct = smuxStruct{
		ServerWrapper: smuxServerWrapper,
		ClientWrapper: smuxClientWrapper,
	}
}

type SmuxConfig struct {
	KeepAliveInterval int
	KeepAliveTimeout  int
	AESKey            string // 如果设置了, 就开启AES加密, 不过还是推荐在上层开启SSL加密为好.
	DisableXOR        bool
}

type smuxServerSideListener struct {
	listener   *smux.Session
	isclose    bool
	aes        *aesStruct
	disableXOR bool
}

func smuxServerWrapper(conn io.ReadWriteCloser, cfg ...SmuxConfig) *smuxServerSideListener {
	scfg := smux.DefaultConfig()
	if len(cfg) != 0 {
		if cfg[0].KeepAliveInterval != 0 {
			scfg.KeepAliveInterval = time.Duration(cfg[0].KeepAliveInterval) * time.Second
		}
		if cfg[0].KeepAliveTimeout != 0 {
			scfg.KeepAliveTimeout = time.Duration(cfg[0].KeepAliveTimeout) * time.Second
		}
	} else {
		scfg.KeepAliveInterval = time.Duration(2) * time.Second
		scfg.KeepAliveTimeout = time.Duration(7) * time.Second
	}

	listener, err := smux.Server(conn, scfg)
	Panicerr(err)
	m := &smuxServerSideListener{listener: listener}

	if len(cfg) != 0 && cfg[0].AESKey != "" {
		m.aes = getAES(cfg[0].AESKey)
	}

	go func() {
		for {
			m.isclose = m.listener.IsClosed()
			if m.isclose {
				break
			}
			sleep(0.1)
		}
	}()

	return m
}

type smuxServerSideConnection struct {
	stream     *smux.Stream
	isclose    bool
	aes        *aesStruct
	disableXOR bool
}

func (m *smuxServerSideListener) accept() chan *smuxServerSideConnection {
	ch := make(chan *smuxServerSideConnection)

	go func() {
		err := Try(func() {
			for {
				stream, err := m.listener.AcceptStream()
				Panicerr(err)

				m := &smuxServerSideConnection{stream: stream, aes: m.aes, disableXOR: m.disableXOR}

				go func() {
					<-m.stream.GetDieCh()
					//lg.trace("Stream is closed")
					m.isclose = true
				}()

				ch <- m
			}
		})
		Lg.Trace("smux接收新连接的时候报错:", err, "session为:", m.listener)
		close(ch)
	}()
	return ch
}

func (m *smuxServerSideConnection) Send(data map[string]string, timeout ...int) {
	m.send(data, timeout...)
}

func (m *smuxServerSideConnection) send(data map[string]string, timeout ...int) {
	if len(timeout) != 0 {
		m.stream.SetWriteDeadline(time.Now().Add(time.Duration(timeout[0]) * time.Second))
	}

	text := map2bin(data)
	if m.aes != nil {
		text = m.aes.Encrypt(text)
	}

	xorkey := make([]byte, 4)
	rand.Read(xorkey)

	if !m.disableXOR {
		_, err := m.stream.Write(xorkey)
		Panicerr(err)
	}

	btlen := make([]byte, 4)
	binary.LittleEndian.PutUint32(btlen, uint32(len(text)))

	if !m.disableXOR {

		_, err := m.stream.Write([]byte(xor(Str(btlen), Str(xorkey))))
		Panicerr(err)
	} else {
		_, err := m.stream.Write(btlen)
		Panicerr(err)
	}

	if !m.disableXOR {
		_, err := m.stream.Write([]byte(xor(text, Str(xorkey))))
		Panicerr(err)
	} else {
		_, err := m.stream.Write([]byte(text))
		Panicerr(err)
	}

	m.stream.SetWriteDeadline(time.Time{})
}

func (m *smuxServerSideConnection) recv(timeout ...int) (data map[string]string) {
	if len(timeout) != 0 {
		m.stream.SetReadDeadline(time.Now().Add(time.Duration(timeout[0]) * time.Second))
	}

	var xorkey string
	if !m.disableXOR {
		headerxorkeylen := 4
		buf := make([]byte, headerxorkeylen)

		for {
			n, err := m.stream.Read(buf)
			if err != nil {
				if err.Error() == "timeout" {
					return nil
				}
				Panicerr(err)
			}

			xorkey = xorkey + string(buf[:n])

			if len(xorkey) != Int(headerxorkeylen) {
				buf = make([]byte, Int(headerxorkeylen)-len(xorkey))
			} else {
				break
			}
		}
	}

	headertotallen := 4
	totalblen := ""
	buf := make([]byte, headertotallen)
	for {
		n, err := m.stream.Read(buf)
		if err != nil {
			if err.Error() == "timeout" {
				return nil
			}
			Panicerr(err)
		}
		totalblen = totalblen + string(buf[:n])

		if len(totalblen) != Int(headertotallen) {
			buf = make([]byte, Int(headertotallen)-len(totalblen))
		} else {
			break
		}
	}
	if !m.disableXOR {

		totalblen = xor(totalblen, xorkey)

	}
	totallen := binary.LittleEndian.Uint32([]byte(totalblen))

	totaldata := ""
	buf = make([]byte, totallen)
	for {
		n, err := m.stream.Read(buf)
		Panicerr(err)

		totaldata = totaldata + string(buf[:n])

		if len(totaldata) != Int(totallen) {
			buf = make([]byte, Int(totallen)-len(totaldata))
		} else {
			break
		}
	}

	if !m.disableXOR {
		totaldata = xor(totaldata, xorkey)
	}
	if m.aes != nil {
		totaldata = m.aes.Decrypt(totaldata)
	}
	data = bin2map(totaldata)

	m.stream.SetReadDeadline(time.Time{})

	return
}

func (m *smuxServerSideConnection) close() {
	if !m.isclose {
		m.stream.Close()
	}
}

type smuxClientSideSession struct {
	session    *smux.Session
	isclose    bool
	aes        *aesStruct
	disableXOR bool
}

func smuxClientWrapper(conn io.ReadWriteCloser, cfg ...SmuxConfig) *smuxClientSideSession {
	scfg := smux.DefaultConfig()
	if len(cfg) != 0 {
		if cfg[0].KeepAliveInterval != 0 {
			scfg.KeepAliveInterval = time.Duration(cfg[0].KeepAliveInterval) * time.Second
		}
		if cfg[0].KeepAliveTimeout != 0 {
			scfg.KeepAliveTimeout = time.Duration(cfg[0].KeepAliveTimeout) * time.Second
		}
	} else {
		scfg.KeepAliveInterval = time.Duration(2) * time.Second
		scfg.KeepAliveTimeout = time.Duration(7) * time.Second
	}

	session, err := smux.Client(conn, scfg)
	Panicerr(err)

	m := &smuxClientSideSession{session: session}
	if len(cfg) != 0 && cfg[0].AESKey != "" {
		m.aes = getAES(cfg[0].AESKey)
	}

	go func() {
		for {
			m.isclose = m.session.IsClosed()
			if m.isclose {
				break
			}
			sleep(0.1)
		}
	}()

	return m
}

type smuxClientSideConnection struct {
	stream     *smux.Stream
	isclose    bool
	aes        *aesStruct
	disableXOR bool
}

func (m *smuxClientSideSession) connect() *smuxClientSideConnection {
	stream, err := m.session.OpenStream()
	Panicerr(err)

	mm := &smuxClientSideConnection{stream: stream, aes: m.aes, disableXOR: m.disableXOR}

	go func() {
		<-mm.stream.GetDieCh()
		mm.isclose = true
	}()

	return mm
}

func (m *smuxClientSideSession) close() {
	if !m.isclose {
		m.isclose = true
		Try(func() {
			m.session.Close()
		})
	}
}

func (m *smuxClientSideConnection) send(data map[string]string, timeout ...int) {
	if len(timeout) != 0 {
		m.stream.SetWriteDeadline(time.Now().Add(time.Duration(timeout[0]) * time.Second))
	}

	text := map2bin(data)
	if m.aes != nil {
		text = m.aes.Encrypt(text)
	}

	xorkey := make([]byte, 4)
	rand.Read(xorkey)

	if !m.disableXOR {
		_, err := m.stream.Write(xorkey)
		Panicerr(err)
	}

	btlen := make([]byte, 4)
	binary.LittleEndian.PutUint32(btlen, uint32(len(text)))

	if !m.disableXOR {

		_, err := m.stream.Write([]byte(xor(Str(btlen), Str(xorkey))))
		Panicerr(err)
	} else {
		_, err := m.stream.Write(btlen)
		Panicerr(err)
	}

	if !m.disableXOR {
		_, err := m.stream.Write([]byte(xor(text, Str(xorkey))))
		Panicerr(err)
	} else {
		_, err := m.stream.Write([]byte(text))
		Panicerr(err)
	}

	m.stream.SetWriteDeadline(time.Time{})
}

func (m *smuxClientSideConnection) recv(timeout ...int) (data map[string]string) {
	if len(timeout) != 0 {
		m.stream.SetReadDeadline(time.Now().Add(time.Duration(timeout[0]) * time.Second))
	}

	var xorkey string
	if !m.disableXOR {
		headerxorkeylen := 4
		buf := make([]byte, headerxorkeylen)

		for {
			n, err := m.stream.Read(buf)
			if err != nil {
				if err.Error() == "timeout" {
					return nil
				}
				Panicerr(err)
			}

			xorkey = xorkey + string(buf[:n])

			if len(xorkey) != Int(headerxorkeylen) {
				buf = make([]byte, Int(headerxorkeylen)-len(xorkey))
			} else {
				break
			}
		}
	}

	headertotallen := 4
	totalblen := ""
	buf := make([]byte, headertotallen)
	for {
		n, err := m.stream.Read(buf)
		if err != nil {
			if err.Error() == "timeout" {
				return nil
			}
			Panicerr(err)
		}
		totalblen = totalblen + string(buf[:n])

		if len(totalblen) != Int(headertotallen) {
			buf = make([]byte, Int(headertotallen)-len(totalblen))
		} else {
			break
		}
	}
	if !m.disableXOR {

		totalblen = xor(totalblen, xorkey)

	}
	totallen := binary.LittleEndian.Uint32([]byte(totalblen))

	totaldata := ""
	buf = make([]byte, totallen)
	for {
		n, err := m.stream.Read(buf)
		Panicerr(err)

		totaldata = totaldata + string(buf[:n])

		if len(totaldata) != Int(totallen) {
			buf = make([]byte, Int(totallen)-len(totaldata))
		} else {
			break
		}
	}

	if !m.disableXOR {
		totaldata = xor(totaldata, xorkey)
	}
	if m.aes != nil {
		totaldata = m.aes.Decrypt(totaldata)
	}
	data = bin2map(totaldata)

	m.stream.SetReadDeadline(time.Time{})

	return
}

func (m *smuxClientSideConnection) close() {
	if !m.isclose {
		m.stream.Close()
	}
}
