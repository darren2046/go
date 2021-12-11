package golanglibs

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"time"

	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
)

type kcpStruct struct {
	Listen     func(host string, port int, key string, salt string) *kcpServerSideListener
	Connect    func(host string, port int, key string, salt string) *kcpClientSideConn
	RawListen  func(host string, port int, key string, salt string) *kcpRawServerSideListener
	RawConnect func(host string, port int, key string, salt string) *kcp.UDPSession
}

var kcpstruct kcpStruct

func init() {
	kcpstruct = kcpStruct{
		Listen:     kcpListen,
		Connect:    kcpConnect,
		RawListen:  kcpRawListen,
		RawConnect: kcpRawConnect,
	}
}

var kcp_ping uint32 = 4000000000
var kcp_pong uint32 = 4000000001
var kcp_close uint32 = 4000000002
var kcp_heartbeat_second int = 20

func kcpRecvSendChanIsClosed(ch chan map[string]string) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

// KCP - Server

type kcpServerSideConn struct {
	conn          *kcp.UDPSession
	isclose       bool
	readtimeout   int
	writetimeout  int
	sendchan      chan map[string]string
	recvchan      chan map[string]string
	heartbeatTime float64
}

type kcpServerSideListener struct {
	listener *kcp.Listener
	isclose  bool
}

func kcpListen(host string, port int, key string, salt string) *kcpServerSideListener {
	block, err := kcp.NewAESBlockCrypt(pbkdf2.Key([]byte(key), []byte(salt), 4096, 32, sha256.New))
	panicerr(err)

	l, err := kcp.ListenWithOptions(host+":"+Str(port), block, 10, 3)
	panicerr(err)

	l.SetDSCP(46)
	l.SetReadBuffer(4194304)
	l.SetWriteBuffer(4194304)

	return &kcpServerSideListener{listener: l}
}

func (m *kcpServerSideListener) accept() chan *kcpServerSideConn {
	ch := make(chan *kcpServerSideConn)

	go func() {
		for {
			var mc *kcpServerSideConn
			if err := try(func() {
				c, err := m.listener.AcceptKCP()
				if err != nil {
					if String("io: read/write on closed pipe").In(err.Error()) || String("use of closed network connection").In(err.Error()) {
						close(ch)
						m.close()
					}
					panicerr(err)
				}
				//
				c.SetNoDelay(0, 20, 2, 1)
				c.SetMtu(1400)
				c.SetWindowSize(1024, 1024)
				c.SetACKNoDelay(false)

				mc = &kcpServerSideConn{
					conn:          c,
					isclose:       false,
					readtimeout:   120,
					writetimeout:  120,
					heartbeatTime: Time.Now(),
					sendchan:      make(chan map[string]string),
					recvchan:      make(chan map[string]string, 10),
				}
			}).Error; err != nil {
				if m.isclose {
					break
				} else {
					sleep(1)
					continue
				}
			}

			// sender
			go func(mc *kcpServerSideConn) {
				try(func() {
					for {
						// 如果连接，关闭，关掉chan，退出
						if mc.isclose {
							if !kcpRecvSendChanIsClosed(mc.sendchan) {
								if !kcpRecvSendChanIsClosed(mc.sendchan) {
									close(mc.sendchan)
								}
							}
							break
						}
						// 如果没有数据就等待1秒再查看，有数据就处理
						select {
						case data, ok := <-mc.sendchan:
							if ok {
								mc.conn.SetWriteDeadline(time.Now().Add(time.Duration(mc.writetimeout) * time.Second))

								btlen := make([]byte, 4)
								// 要发送的字符串
								text := map2bin(data)
								binary.LittleEndian.PutUint32(btlen, uint32(len(text)))

								_, err := mc.conn.Write(btlen)
								if err != nil {
									if mc.isclose {
										break
									} else {
										panicerr(err)
									}
								}

								_, err = mc.conn.Write([]byte(text))
								if err != nil {
									if mc.isclose {
										break
									} else {
										panicerr(err)
									}
								}
							}
						case <-time.After(getTimeDuration(1)):

						}
					}
				}).except(func(e error) {
				})
			}(mc)

			// receiver
			go func(mc *kcpServerSideConn) {
				try(func() {
					for {
						mc.conn.SetReadDeadline(time.Now().Add(time.Duration(mc.readtimeout) * time.Second))

						// 读取总长度, 然后读取所有数据
						totalblen := make([]byte, 4096)
						n, err := mc.conn.Read(totalblen)
						if err != nil {
							if mc.isclose || String("io: read/write on closed pipe").In(err.Error()) {
								if mc.isclose {
								} else {
								}
								mc.isclose = true
								if !kcpRecvSendChanIsClosed(mc.recvchan) {
									close(mc.recvchan)
								}
								break
							} else {
								panicerr(err)
							}
						}
						if n != 4 {
							continue
						}

						totallen := binary.LittleEndian.Uint32(totalblen)

						// 如果是心跳，更新心跳时间，回复它，处理下一个包
						if totallen == kcp_ping {
							btlen := make([]byte, 4)
							binary.LittleEndian.PutUint32(btlen, kcp_pong)
							_, err := mc.conn.Write(btlen)
							if err != nil {
								if mc.isclose {
									break
								} else {
									panicerr(err)
								}
							}
							mc.heartbeatTime = Time.Now()
							continue
						}

						totaldata := ""
						buf := make([]byte, 4096)
						for len(totaldata) != Int(totallen) {
							n, err := mc.conn.Read(buf)
							if err != nil {
								if mc.isclose {
									if !kcpRecvSendChanIsClosed(mc.recvchan) {
										close(mc.recvchan)
									}
									break
								} else {
									panicerr(err)
								}
							}

							if n == 4 {
								buflen := binary.LittleEndian.Uint32(buf)
								// 如果是心跳，更新心跳时间，回复它，处理下一个包
								if buflen == kcp_ping {
									btlen := make([]byte, 4)
									binary.LittleEndian.PutUint32(btlen, kcp_pong)
									_, err := mc.conn.Write(btlen)
									if err != nil {
										if mc.isclose {
											break
										} else {
											panicerr(err)
										}
									}
									mc.heartbeatTime = Time.Now()
									continue
								}
							}

							totaldata = totaldata + string(buf[:n])
						}

						res := bin2map(totaldata)
						mc.recvchan <- res
					}
				}).except(func(e error) {
				})
			}(mc)

			// heartbeat checker
			// 如果连接被主动关闭，退出，如果3次没有收到心跳，关闭连接
			go func(mc *kcpServerSideConn) {
				try(func() {
					for {
						if mc.isclose {
							break
						}
						if Time.Now()-mc.heartbeatTime > Float64(kcp_heartbeat_second)*3 {
							if !kcpRecvSendChanIsClosed(mc.recvchan) {
								close(mc.recvchan)
							}
							if !kcpRecvSendChanIsClosed(mc.sendchan) {
								close(mc.sendchan)
							}
							mc.close()
							break
						}
						sleep(1)
					}
				}).except(func(e error) {
				})
			}(mc)

			ch <- mc
		}
	}()

	return ch
}

func (m *kcpServerSideListener) close() {
	if !m.isclose {
		m.isclose = true
		m.listener.Close()
	}
}

func (m *kcpServerSideConn) close() {
	if !m.isclose {
		m.isclose = true
		try(func() {
			m.conn.Close()
		})
	}
	if !kcpRecvSendChanIsClosed(m.sendchan) {
		close(m.sendchan)
	}
	if !kcpRecvSendChanIsClosed(m.recvchan) {
		close(m.recvchan)
	}
}

func (m *kcpServerSideConn) send(data map[string]string) {
	if m.isclose {
		err := errors.New("连接已关闭，不可发送数据")
		panicerr(err)
	}
	m.sendchan <- data
}

func (m *kcpServerSideConn) recv(timeoutSecond ...int) (res map[string]string) {
	if m.isclose {
		err := errors.New("连接已关闭，不可收取数据")
		panicerr(err)
	}
	if len(timeoutSecond) == 0 {
		res = <-m.recvchan
	} else {
		for range Range(timeoutSecond[0] * 10) {
			select {
			case data, ok := <-m.recvchan:
				if ok {
					res = data
					break
				}
			default:
				sleep(0.1)
			}
			if res != nil {
				break
			}
		}
	}

	if res == nil {
		if m.isclose {
			err := errors.New("连接已关闭，不可收取数据")
			panicerr(err)
		}
	}
	return
}

// KCP - Client

type kcpClientSideConn struct {
	conn          *kcp.UDPSession
	isclose       bool
	readtimeout   int
	writetimeout  int
	recvchan      chan map[string]string
	sendchan      chan map[string]string
	heartbeatTime float64
}

func kcpConnect(host string, port int, key string, salt string) *kcpClientSideConn {
	block, err := kcp.NewAESBlockCrypt(pbkdf2.Key([]byte(key), []byte(salt), 4096, 32, sha256.New))
	panicerr(err)
	conn, err := kcp.DialWithOptions(host+":"+Str(port), block, 10, 3)
	panicerr(err)

	conn.SetMtu(1400)
	conn.SetWriteDelay(false)
	conn.SetNoDelay(0, 20, 2, 1)
	conn.SetWindowSize(128, 1024)
	conn.SetACKNoDelay(false)
	conn.SetDSCP(46)
	conn.SetReadBuffer(4194304)
	conn.SetWriteBuffer(4194304)

	m := &kcpClientSideConn{
		conn:          conn,
		isclose:       false,
		readtimeout:   120,
		writetimeout:  120,
		sendchan:      make(chan map[string]string),
		recvchan:      make(chan map[string]string, 10),
		heartbeatTime: Time.Now(),
	}

	// sender
	go func(m *kcpClientSideConn) {
		try(func() {
			for {
				// 如果连接，关闭，关掉chan，退出
				if m.isclose {
					if !kcpRecvSendChanIsClosed(m.sendchan) {
						close(m.sendchan)
					}
					break
				}
				// 如果没有数据就等待0.3秒再查看，有数据就处理
				select {
				case data, ok := <-m.sendchan:
					if ok {
						m.conn.SetWriteDeadline(time.Now().Add(time.Duration(m.writetimeout) * time.Second))
						// 要发送的字符串
						text := map2bin(data)
						// 4字节的二进制无符号int
						btlen := make([]byte, 4)

						binary.LittleEndian.PutUint32(btlen, uint32(len(text)))

						_, err := m.conn.Write(btlen)
						if err != nil {
							if m.isclose {
								break
							} else {
								panicerr(err)
							}
						}

						_, err = m.conn.Write([]byte(text))
						if err != nil {
							if m.isclose {
								break
							} else {
								panicerr(err)
							}
						}
					}
				case <-time.After(getTimeDuration(1)):
				}
			}
		}).except(func(e error) {
		})
	}(m)

	// receiver
	go func(m *kcpClientSideConn) {
		try(func() {
			for {
				m.conn.SetReadDeadline(time.Now().Add(time.Duration(m.readtimeout) * time.Second))

				// 读取总长度, 然后读取所有数据
				totalblen := make([]byte, 4096)
				n, err := m.conn.Read(totalblen)
				if err != nil {
					if m.isclose {
						if !kcpRecvSendChanIsClosed(m.recvchan) {
							close(m.recvchan)
						}
						break
					} else {
						panicerr(err)
					}
				}
				if n != 4 {
					continue
				}

				totallen := binary.LittleEndian.Uint32(totalblen)

				// 如果是心跳的回应, 继续下一个包
				if totallen == kcp_pong {
					m.heartbeatTime = Time.Now()
					continue
				}

				totaldata := ""
				buf := make([]byte, 4096)
				for len(totaldata) != Int(totallen) {
					n, err := m.conn.Read(buf)
					if err != nil {
						if m.isclose {
							if !kcpRecvSendChanIsClosed(m.recvchan) {
								close(m.recvchan)
							}
							break
						} else {
							panicerr(err)
						}
					}

					if n == 4 {
						buflen := binary.LittleEndian.Uint32(buf)
						// 如果是心跳的回应, 继续下一个包
						if buflen == kcp_pong {
							m.heartbeatTime = Time.Now()
							continue
						}
					}

					totaldata = totaldata + string(buf[:n])
				}

				res := bin2map(totaldata)
				m.recvchan <- res
			}
		}).except(func(e error) {
		})
	}(m)

	// heartbeat
	go func(m *kcpClientSideConn) {
		try(func() {
			for {
				btlen := make([]byte, 4)
				binary.LittleEndian.PutUint32(btlen, kcp_ping)
				_, err := m.conn.Write(btlen)
				if err != nil {
					if m.isclose {
						if !kcpRecvSendChanIsClosed(m.sendchan) {
							close(m.sendchan)
						}
						if !kcpRecvSendChanIsClosed(m.recvchan) {
							close(m.recvchan)
						}
						break
					} else {
						panicerr(err)
					}
				}
				if Time.Now()-m.heartbeatTime > Float64(kcp_heartbeat_second)*3 {
					m.close()
				}
				sleep(kcp_heartbeat_second)
			}
		}).except(func(e error) {
		})
	}(m)

	return m
}

func (m *kcpClientSideConn) send(data map[string]string) {
	if m.isclose {
		err := errors.New("连接已关闭，不可发送数据")
		panicerr(err)
	}
	m.sendchan <- data
}

func (m *kcpClientSideConn) recv(timeoutSecond ...int) (res map[string]string) {
	if m.isclose {
		err := errors.New("连接已关闭，不可收取数据")
		panicerr(err)
	}
	if len(timeoutSecond) == 0 {
		res = <-m.recvchan
	} else {
		for range Range(timeoutSecond[0] * 10) {
			select {
			case data, ok := <-m.recvchan:
				if ok {
					res = data
					break
				}
			default:
				sleep(0.1)
			}
			if res != nil {
				break
			}
		}
	}

	if res == nil {
		if m.isclose {
			err := errors.New("连接已关闭，不可收取数据")
			panicerr(err)
		}
	}
	return
}

func (m *kcpClientSideConn) close() {
	if !m.isclose {
		m.isclose = true
		m.conn.Close()
	}

	if !kcpRecvSendChanIsClosed(m.sendchan) {
		close(m.sendchan)
	}
	if !kcpRecvSendChanIsClosed(m.recvchan) {
		close(m.recvchan)
	}
}
