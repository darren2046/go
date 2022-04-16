package golanglibs

import (
	"testing"
)

func TestTelegram(t *testing.T) {
	ident := Open("telegram.app.hash.txt").Read() // ident: 12334566,72aae231fa287ec6b04c39a27ae94ed0
	AppID := Int32(ident.Split(",")[0].Strip().S)
	AppHash := ident.Split(",")[1].Strip().S
	Lg.Trace(AppID, AppHash)
	tg := Tools.Telegram(
		AppID,
		AppHash,
		TelegramConfig{
			SessionFile: "telegram.session.file.json",
		},
	)

	// Lg.Trace(tg)
	// Lg.Trace()

	// 列出chat
	// for _, chat := range tg.Chats() {
	// 	Lg.Trace(chat.Type, chat.Title, chat.Username)
	// 	if chat.Username == "zzzbot" {
	// 		chat.Send("test")
	// 	}
	// 	// if chat.Title == "test" {
	// 	// 	for _, msg := range chat.History(3) {
	// 	// 		Lg.Debug(msg)
	// 	// 	}
	// 	// }
	// }

	// 根据username获取历史记录
	// p := tg.ResolvePeerByUsername("12345678910")
	// Lg.Debug(p)
	// Lg.Debug(p.History(10))

	// 根据id和accesshash获取历史记录
	p := tg.NewInputPeer("channel", 1234567890, -987654321)
	Lg.Trace(len(p.History(1)))
}
