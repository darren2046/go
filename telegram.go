package golanglibs

import (
	"fmt"
	stdlog "log"

	"encoding/json"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/3bl3gamer/tgclient"
	"github.com/3bl3gamer/tgclient/mtproto"
	"github.com/ansel1/merry"
	"github.com/fatih/color"
	"golang.org/x/net/proxy"
)

// ------- Telegam Log ----------

type LogHandler struct {
	mtproto.ColorLogHandler
	ConsoleMaxLevel mtproto.LogLevel
	ErrorFileLoger  *stdlog.Logger
	DebugFileLoger  *stdlog.Logger
	ConsoleLogger   *stdlog.Logger
}

func (h LogHandler) Log(level mtproto.LogLevel, err error, msg string, args ...interface{}) {
	text := h.StringifyLog(level, err, msg, args...)
	text = h.AddLevelPrevix(level, text)
	if level <= h.ConsoleMaxLevel {
		h.ConsoleLogger.Print(h.AddLevelColor(level, text))
	}
	if level <= mtproto.ERROR {
		h.ErrorFileLoger.Print(text)
	}
	h.DebugFileLoger.Print(text)
}

func (h LogHandler) Message(isIncoming bool, msg mtproto.TL, id int64) {
	h.Log(mtproto.DEBUG, nil, h.StringifyMessage(isIncoming, msg, id))
}

// ------- Telegram Login ----------

type TelegramConfig struct {
	SessionFile string
	Proxy       string // socks5://user:pass@host:port
}

type TelegramStruct struct {
	tg *tgclient.TGClient
}

func getTelegram(AppID int32, AppHash string, config ...TelegramConfig) *TelegramStruct {
	cfg := &mtproto.AppConfig{
		AppID:          AppID,
		AppHash:        AppHash,
		AppVersion:     "0.0.1",
		DeviceModel:    "Unknown",
		SystemVersion:  runtime.GOOS + "/" + runtime.GOARCH,
		SystemLangCode: "en",
		LangPack:       "",
		LangCode:       "en",
	}

	var sessStore *mtproto.SessFileStore
	var dialer proxy.Dialer
	if len(config) != 0 {
		cfg := config[0]
		if cfg.SessionFile != "" {
			sessStore = &mtproto.SessFileStore{FPath: cfg.SessionFile}
		} else {
			sessStore = &mtproto.SessFileStore{FPath: "/dev/null"}
		}

		if cfg.Proxy != "" {
			u := Tools.URL(cfg.Proxy).Parse()
			if u.Schema != "socks5" {
				Panicerr("不支持的代理协议: " + u.Schema)
			}

			var auth *proxy.Auth
			if u.User != "" || u.Pass != "" {
				auth = &proxy.Auth{
					User:     u.User,
					Password: u.Pass,
				}
			}
			var err error
			dialer, err = proxy.SOCKS5("tcp", u.Host+":"+u.Port, auth, proxy.Direct)
			Panicerr(err)
		}
	}
	// Log
	tgLogHandler := LogHandler{
		ConsoleMaxLevel: mtproto.INFO,
		DebugFileLoger:  stdlog.New(Open("/dev/null", "w").fd, "", stdlog.LstdFlags),
		ErrorFileLoger:  stdlog.New(Open("/dev/null", "w").fd, "", stdlog.LstdFlags),
		ConsoleLogger:   stdlog.New(color.Error, "", stdlog.LstdFlags),
	}

	tg := tgclient.NewTGClientExt(cfg, sessStore, tgLogHandler, dialer)

	err := tg.InitAndConnect()
	Panicerr(err)

	res, err := tg.AuthExt(mtproto.ScanfAuthDataProvider{}, mtproto.TL_users_getUsers{ID: []mtproto.TL{mtproto.TL_inputUserSelf{}}})
	Panicerr(err)

	users, ok := res.(mtproto.VectorObject)
	if !ok {
		Panicerr("获取不到自己的用户信息")
	}

	me := users[0].(mtproto.TL_user)
	Lg.Trace("以身份登录telegram:", me)

	return &TelegramStruct{
		tg: tg,
	}
}

// --------- Chat ---------

type tgChatType int8

const (
	tgChatUser tgChatType = iota
	tgChatGroup
	tgChatChannel
)

func (t tgChatType) String() string {
	switch t {
	case tgChatUser:
		return "user"
	case tgChatGroup:
		return "group"
	case tgChatChannel:
		return "channel"
	default:
		Panicerr("未知的chat类型")
		return ""
	}
}

func (t *tgChatType) UnmarshalJSON(buf []byte) {
	var s string
	err := json.Unmarshal(buf, &s)
	Panicerr(err)

	switch s {
	case "user":
		*t = tgChatUser
	case "group":
		*t = tgChatGroup
	case "channel":
		*t = tgChatChannel
	default:
		Panicerr("未知的chat类型")
	}
}

type TelegramChatStruct struct {
	ID            int64
	Title         string
	Username      string
	LastMessageID int32
	Type          tgChatType
	Obj           mtproto.TL
	tg            *tgclient.TGClient
}

func tgExtractDialogsData(tg *tgclient.TGClient, dialogs []mtproto.TL, chats []mtproto.TL, users []mtproto.TL) ([]*TelegramChatStruct, error) {
	chatsByID := make(map[int64]mtproto.TL_chat)
	channelsByID := make(map[int64]mtproto.TL_channel)
	for _, chatTL := range chats {
		switch chat := chatTL.(type) {
		case mtproto.TL_chat:
			chatsByID[chat.ID] = chat
		case mtproto.TL_chatForbidden:
			chatsByID[chat.ID] = mtproto.TL_chat{ID: chat.ID, Title: chat.Title}
		case mtproto.TL_channel:
			channelsByID[chat.ID] = chat
		case mtproto.TL_channelForbidden:
			channelsByID[chat.ID] = mtproto.TL_channel{ID: chat.ID, Title: chat.Title, AccessHash: chat.AccessHash, Megagroup: chat.Megagroup}
		default:
			return nil, merry.Wrap(mtproto.WrongRespError(chatTL))
		}
	}
	usersByID := make(map[int64]mtproto.TL_user)
	for _, userTL := range users {
		user := userTL.(mtproto.TL_user)
		usersByID[user.ID] = user
	}
	extractedChats := make([]*TelegramChatStruct, len(dialogs))
	for i, chatTL := range dialogs {
		dialog := chatTL.(mtproto.TL_dialog)
		ext := &TelegramChatStruct{LastMessageID: dialog.TopMessage, tg: tg}
		switch peer := dialog.Peer.(type) {
		case mtproto.TL_peerUser:
			user := usersByID[peer.UserID]
			ext.ID = user.ID
			ext.Title = strings.TrimSpace(user.FirstName + " " + user.LastName)
			ext.Username = user.Username
			ext.Type = tgChatUser
			ext.Obj = user
		case mtproto.TL_peerChat:
			chat := chatsByID[peer.ChatID]
			ext.ID = chat.ID
			ext.Title = chat.Title
			ext.Type = tgChatGroup
			ext.Obj = chat
		case mtproto.TL_peerChannel:
			channel := channelsByID[peer.ChannelID]
			ext.ID = channel.ID
			ext.Title = channel.Title
			ext.Username = channel.Username
			ext.Type = tgChatChannel
			if channel.Megagroup {
				ext.Type = tgChatGroup
			}
			ext.Obj = channel
		default:
			return nil, merry.Wrap(mtproto.WrongRespError(dialog.Peer))
		}
		extractedChats[i] = ext
	}
	return extractedChats, nil
}

func tgGetMessageStamp(msgTL mtproto.TL) (int32, error) {
	switch msg := msgTL.(type) {
	case mtproto.TL_message:
		return msg.Date, nil
	case mtproto.TL_messageService:
		return msg.Date, nil
	default:
		return 0, merry.Wrap(mtproto.WrongRespError(msg))
	}
}

func tgLoadChats(tg *tgclient.TGClient) []*TelegramChatStruct {
	chats := make([]*TelegramChatStruct, 0)
	offsetDate := int32(0)
	for {
		res := tg.SendSyncRetry(mtproto.TL_messages_getDialogs{
			OffsetPeer: mtproto.TL_inputPeerEmpty{},
			OffsetDate: offsetDate,
			Limit:      100,
		}, time.Second, 0, 30*time.Second)

		switch slice := res.(type) {
		case mtproto.TL_messages_dialogs:
			chats, err := tgExtractDialogsData(tg, slice.Dialogs, slice.Chats, slice.Users)
			Panicerr(err)

			return chats
		case mtproto.TL_messages_dialogsSlice:
			group, err := tgExtractDialogsData(tg, slice.Dialogs, slice.Chats, slice.Users)
			Panicerr(err)

			for _, d := range group {
				chats = append(chats, d) //TODO: check duplicates
			}

			offsetDate, err = tgGetMessageStamp(slice.Messages[len(slice.Messages)-1])
			Panicerr(err)

			if len(chats) == int(slice.Count) {
				return chats
			}
			if len(slice.Dialogs) < 100 {
				// fmt.Printf("some chats seem missing: got %d in the end, expected %d; retrying from start\n", len(chats), slice.Count)
				offsetDate = 0
			}
		default:
			Panicerr(merry.Wrap(mtproto.WrongRespError(res)))
		}
	}
}

func (m *TelegramStruct) Chats() []*TelegramChatStruct {
	return tgLoadChats(m.tg)
}

// ---------------- Message History -------------

type tgUserData struct {
	ID                    int64
	Username              string
	FirstName, LastName   string
	PhoneNumber, LangCode string
	IsBot                 bool
}

func (u *tgUserData) Equals(other *mtproto.TL_user) bool {
	// Sometimes Username becomes blank and then becomes filled again.
	// This will produce unnesessary updates in users file. So just ignoring that change.
	return (other.Username == "" || u.Username == other.Username) &&
		u.FirstName == other.FirstName && u.LastName == other.LastName &&
		u.PhoneNumber == other.Phone &&
		u.LangCode == other.LangCode
}

type tgChatData struct {
	ID        int64
	Username  string
	Title     string
	IsChannel bool
}

func (c *tgChatData) Equals(other *tgChatData) bool {
	return c.Username == other.Username && c.Title == other.Title
}

func tgObjToMap(obj mtproto.TL) map[string]interface{} {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	typ := v.Type()
	res := make(map[string]interface{})
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		var val interface{}
		switch value := v.Field(i).Interface().(type) {
		case int64:
			val = strconv.FormatInt(value, 10)
		case mtproto.TL:
			val = tgObjToMap(value)
		case []mtproto.TL:
			vals := make([]interface{}, len(value))
			for i, item := range value {
				vals[i] = tgObjToMap(item)
			}
			val = vals
		default:
			val = value
		}
		res[field.Name] = val
	}
	res["_"] = typ.Name()
	return res
}

// 未区分转发，未判断多媒体，未标明是否被编辑和是否是回复某个消息
type tgMessageStruct struct {
	Chat        *tgChatData
	Message     string
	User        *tgUserData
	Time        int64
	ReplyMarkup map[string]interface{}
	Entities    []interface{}
}

func (m *TelegramChatStruct) History(limit int32) (resmsgs []*tgMessageStruct) {
	params := mtproto.TL_messages_getHistory{
		Peer:  m.GetPeer(),
		Limit: limit,
	}

	res := m.tg.SendSyncRetry(params, time.Second, 0, 30*time.Second)

	var msgs []mtproto.TL
	var musers []mtproto.TL
	var mchats []mtproto.TL

	switch messages := res.(type) {
	case mtproto.TL_messages_messages:
		msgs = messages.Messages
		musers = messages.Users
		mchats = messages.Chats
	case mtproto.TL_messages_messagesSlice:
		msgs = messages.Messages
		musers = messages.Users
		mchats = messages.Chats
	case mtproto.TL_messages_channelMessages:
		msgs = messages.Messages
		musers = messages.Users
		mchats = messages.Chats
	default:
		panic("未知的消息类型")
	}

	tguserdatas := make(map[int64]*tgUserData)

	for _, userTL := range musers {
		tgUser, ok := userTL.(mtproto.TL_user)
		if !ok {
			panic("未知的tgUserTL")
		}
		tguserdatas[tgUser.ID] = &tgUserData{
			ID:          tgUser.ID,
			FirstName:   tgUser.FirstName,
			LastName:    tgUser.LastName,
			Username:    tgUser.Username,
			PhoneNumber: tgUser.Phone,
			LangCode:    tgUser.LangCode,
			IsBot:       tgUser.Bot,
		}
		// Lg.Debug("====> User:", newUser)
	}

	tgchatdatas := make(map[int64]*tgChatData)
	for _, chatTL := range mchats {
		var newChat *tgChatData
		switch c := chatTL.(type) {
		case mtproto.TL_chat:
			newChat = &tgChatData{ID: c.ID, Title: c.Title}
		case mtproto.TL_chatForbidden:
			newChat = &tgChatData{ID: c.ID, Title: c.Title}
		case mtproto.TL_channel:
			newChat = &tgChatData{ID: c.ID, Title: c.Title, Username: c.Username, IsChannel: !c.Megagroup}
		case mtproto.TL_channelForbidden:
			newChat = &tgChatData{ID: c.ID, Title: c.Title, IsChannel: !c.Megagroup}
		default:
			panic("未知的chat类型")
		}

		tgchatdatas[newChat.ID] = newChat

		//Lg.Debug("====>Chat: ", newChat)
	}

	for i := len(msgs) - 1; i >= 0; i-- {
		msg := msgs[i]
		msgMap := tgObjToMap(msg)
		msgMap["_TL_LAYER"] = mtproto.TL_Layer

		Lg.Debug("====>Msg:", msg)

		if Str(msgMap["Message"]) == "" {
			continue
		}
		resmsgs = append(resmsgs, &tgMessageStruct{
			Message:     Str(msgMap["Message"]),
			Chat:        tgchatdatas[Int64(StringMap(msgMap["PeerID"])["ChatID"])],
			User:        tguserdatas[Int64(StringMap(msgMap["FromID"])["UserID"])],
			Time:        Int64(msgMap["Date"]),
			ReplyMarkup: StringMap(msgMap["ReplyMarkup"]),
			Entities:    InterfaceArray(msgMap["Entities"]),
		})
	}

	return
}

func (m *TelegramChatStruct) GetPeer() (inputPeer mtproto.TL) {
	switch peer := m.Obj.(type) {
	case mtproto.TL_user:
		inputPeer = mtproto.TL_inputPeerUser{UserID: peer.ID, AccessHash: peer.AccessHash}
	case mtproto.TL_chat:
		inputPeer = mtproto.TL_inputPeerChat{ChatID: peer.ID}
	case mtproto.TL_channel:
		inputPeer = mtproto.TL_inputPeerChannel{ChannelID: peer.ID, AccessHash: peer.AccessHash}
	default:
		panic("未知的mtproto.TL类型")
	}
	return
}

// --------- send message ---------

func (m *TelegramChatStruct) Send(text string) {
	params := mtproto.TL_messages_sendMessage{
		RandomID: Int64(Time.Now()),
		Peer:     m.GetPeer(),
		Message:  text,
	}

	res := m.tg.SendSyncRetry(params, time.Second, 0, 30*time.Second)

	fmt.Println(Repr(res).S)

	if !Repr(res).Has(text) {
		Panicerr("可能发送消息失败？")
	}
}
