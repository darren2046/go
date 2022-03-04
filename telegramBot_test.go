package golanglibs

import (
	"testing"
)

func TestTelegramBotSendVideo(t *testing.T) {
	tg := Tools.
		TelegramBot("key").
		SetChatID(123456)

	tg.SendVideo("/path/to/file.mp4")
}
