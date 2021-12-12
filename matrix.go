package golanglibs

import (
	"github.com/matrix-org/gomatrix"
)

type matrixStruct struct {
	cli    *gomatrix.Client
	roomID string
}

func getMatrix(homeserverURL string) *matrixStruct {
	cli, err := gomatrix.NewClient(homeserverURL, "", "")
	panicerr(err)

	return &matrixStruct{cli: cli}
}

func (c *matrixStruct) Login(username string, password string) string {
	resp, err := c.cli.Login(&gomatrix.ReqLogin{
		Type:     "m.login.password",
		User:     username,
		Password: password,
	})
	panicerr(err)

	c.SetToken(resp.UserID, resp.AccessToken)

	return resp.AccessToken
}

func (c *matrixStruct) SetToken(userID string, token string) *matrixStruct {
	c.cli.SetCredentials(userID, token)
	return c
}

func (c *matrixStruct) SetRoomID(roomID string) *matrixStruct {
	c.roomID = roomID
	return c
}

func (c *matrixStruct) Send(text string) {
	_, err := c.cli.SendText(c.roomID, text)
	panicerr(err)
}
