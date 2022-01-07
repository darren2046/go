package golanglibs

import (
	"github.com/matrix-org/gomatrix"
)

type MatrixStruct struct {
	cli    *gomatrix.Client
	roomID string
}

func getMatrix(homeserverURL string) *MatrixStruct {
	cli, err := gomatrix.NewClient(homeserverURL, "", "")
	Panicerr(err)

	return &MatrixStruct{cli: cli}
}

func (c *MatrixStruct) Login(username string, password string) string {
	resp, err := c.cli.Login(&gomatrix.ReqLogin{
		Type:     "m.login.password",
		User:     username,
		Password: password,
	})
	Panicerr(err)

	c.SetToken(resp.UserID, resp.AccessToken)

	return resp.AccessToken
}

func (c *MatrixStruct) SetToken(userID string, token string) *MatrixStruct {
	c.cli.SetCredentials(userID, token)
	return c
}

func (c *MatrixStruct) SetRoomID(roomID string) *MatrixStruct {
	c.roomID = roomID
	return c
}

func (c *MatrixStruct) Send(text string) {
	_, err := c.cli.SendText(c.roomID, text)
	Panicerr(err)
}
