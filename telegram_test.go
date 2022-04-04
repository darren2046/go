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

	p := tg.ResolvePeerByUsername("zfzzzshou")
	// Lg.Debug(p)
	Lg.Debug(p.History(10))
}
