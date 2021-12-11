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

func (c *matrixStruct) login(username string, password string) string {
	resp, err := c.cli.Login(&gomatrix.ReqLogin{
		Type:     "m.login.password",
		User:     username,
		Password: password,
	})
	panicerr(err)

	c.setToken(resp.UserID, resp.AccessToken)

	return resp.AccessToken
}

func (c *matrixStruct) setToken(userID string, token string) *matrixStruct {
	c.cli.SetCredentials(userID, token)
	return c
}

func (c *matrixStruct) setRoomID(roomID string) *matrixStruct {
	c.roomID = roomID
	return c
}

func (c *matrixStruct) send(text string) {
	_, err := c.cli.SendText(c.roomID, text)
	panicerr(err)
}
