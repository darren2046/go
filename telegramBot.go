package golanglibs

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type telegramBotStruct struct {
	tg     *tgbotapi.BotAPI
	chatid int64
}

func getTelegramBot(token string) *telegramBotStruct {
	for {
		//lg.trace("开始初始化")
		bot, err := tgbotapi.NewBotAPI(token)
		if err == nil {
			//lg.trace("初始化成功")
			return &telegramBotStruct{tg: bot}
		} else {
			//lg.trace("初始化出错:", err)
			sleep(3)
		}
	}
}

func (m *telegramBotStruct) SetChatID(chatid int64) *telegramBotStruct {
	m.chatid = chatid
	return m
}

func (m *telegramBotStruct) SendFile(path string) tgbotapi.Message {
	var err error
	var msg tgbotapi.Message
	sleepCount := 10
	for {
		msg, err = m.tg.Send(tgbotapi.NewDocumentUpload(m.chatid, path))
		if err == nil {
			break
		}

		sleep(sleepCount)
		sleepCount = sleepCount * 2
	}
	return msg
}

func (m *telegramBotStruct) SendImage(path string) tgbotapi.Message {
	var err error
	var msg tgbotapi.Message
	sleepCount := 10
	for {
		msg, err = m.tg.Send(tgbotapi.NewPhotoUpload(m.chatid, path))
		if err == nil {
			break
		}
		sleep(sleepCount)
		sleepCount = sleepCount * 2
	}
	return msg
}

func (m *telegramBotStruct) SendVideo(path string) tgbotapi.Message {
	var err error
	var msg tgbotapi.Message
	sleepCount := 10
	for {
		msg, err = m.tg.Send(tgbotapi.NewVideoUpload(m.chatid, path))
		if err == nil {
			break
		}
		sleep(sleepCount)
		sleepCount = sleepCount * 2
	}
	return msg
}

func (m *telegramBotStruct) SendAudio(path string) tgbotapi.Message {
	var err error
	var msg tgbotapi.Message
	sleepCount := 10
	for {
		msg, err = m.tg.Send(tgbotapi.NewAudioUpload(m.chatid, path))
		if err == nil {
			break
		}
		sleep(sleepCount)
		sleepCount = sleepCount * 2
	}
	return msg
}

type tgMsgConfig struct {
	parseMode             string
	DisableWebPagePreview bool
	DisableRetryOnError   bool
}

func (m *telegramBotStruct) Send(text string, cfg ...tgMsgConfig) tgbotapi.Message {
	var err error
	var msg tgbotapi.Message

	mm := tgbotapi.NewMessage(m.chatid, text)
	if len(cfg) != 0 {
		mm.ParseMode = cfg[0].parseMode
		mm.DisableWebPagePreview = cfg[0].DisableWebPagePreview
	}

	sleepCount := 10
	for {
		msg, err = m.tg.Send(mm)
		// lg.trace(err)
		if err == nil {
			break
		} else {
			if len(cfg) != 0 && cfg[0].DisableRetryOnError {
				break
			}
			sleep(sleepCount)
			sleepCount = sleepCount * 2
		}
	}
	return msg
}
