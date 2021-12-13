package golanglibs

import (
	"time"

	"github.com/sacOO7/gowebsocket"
)

type websocketStruct struct {
	socket      gowebsocket.Socket
	recvMsgChan chan string
}

func getWebSocket(url string) *websocketStruct {
	socket := gowebsocket.New(url)

	recvMsgChan := make(chan string)

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		recvMsgChan <- message
	}

	socket.Connect()

	return &websocketStruct{
		socket:      socket,
		recvMsgChan: recvMsgChan,
	}
}

func (c *websocketStruct) Send(text string) {
	c.socket.SendText(text)
}

func (c *websocketStruct) Recv(timeout ...int) string {
	if len(timeout) != 0 {
		select {
		case resp := <-c.recvMsgChan:
			return resp
		case <-time.After(getTimeDuration(timeout[0])):
			Panicerr("Timeout while recving data")
		}
	} else {
		return <-c.recvMsgChan
	}
	return ""
}

func (c *websocketStruct) Close() {
	c.socket.Close()
}
