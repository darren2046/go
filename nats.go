package golanglibs

import (
	"github.com/nats-io/nats.go"
)

type natsStruct struct {
	conn *nats.Conn
}

func getNats(server string) *natsStruct {
	// server = "nat://nats.nats.svc.cluster.local"
	conn, err := nats.Connect(
		server,
		nats.DisconnectErrHandler(func(c *nats.Conn, err error) {
			if err != nil {
				Lg.Trace("Disconnected due to: " + err.Error() + ", trying to reconnect")
			} else {
				Lg.Trace("Disconnected normally")
			}
		}),
		nats.ReconnectHandler(func(c *nats.Conn) {
			Lg.Trace("Reconnected [" + c.ConnectedUrl() + "]")
		}),
		nats.ClosedHandler(func(c *nats.Conn) {
			if err := c.LastError(); err != nil {
				Lg.Trace("Connection closed due to: " + c.LastError().Error())
			} else {
				Lg.Trace("Connection closed normally")
			}
		}),
	)
	Panicerr(err)
	return &natsStruct{
		conn: conn,
	}
}

type subjectNatsStruct struct {
	conn    *nats.Conn
	subject string
}

func (m *natsStruct) Subject(subject string) *subjectNatsStruct {
	return &subjectNatsStruct{
		conn:    m.conn,
		subject: subject,
	}
}

func (m *subjectNatsStruct) Publish(message string) {
	err := m.conn.Publish(m.subject, []byte(message))
	Panicerr(err)
}

func (m *subjectNatsStruct) Subscribe() chan string {
	subscribeChan := make(chan string)
	go func() {
		_, err := m.conn.Subscribe(m.subject, func(msg *nats.Msg) {
			subscribeChan <- string(msg.Data)
		})
		if err != nil {
			Lg.Trace("Error while subscribe subject \"" + m.subject + "\":" + err.Error())
			close(subscribeChan)
		}
	}()
	return subscribeChan
}

func (m *subjectNatsStruct) Flush() {
	err := m.conn.Flush()
	Panicerr(err)
}
