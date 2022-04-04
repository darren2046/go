package golanglibs

import (
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
	AccessHash            int64
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
type TelegramMessageStruct struct {
	Chat        *tgChatData
	Message     string
	User        *tgUserData
	Time        int64
	ReplyMarkup map[string]interface{}
	Entities    []interface{}
	ID          int64
	File        *TelegramFileInfo
	Action      string
}

// 不设置offset则获取最新的limit条数据
// 设置了offset的话，例如如果设置offset的id是4，那就是从id小于4的message开始往上获取limit个数据。（不包括id为4的message）
func (m *TelegramChatStruct) History(limit int32, offset ...int32) (resmsgs []*TelegramMessageStruct) {
	return historyMessage(m.tg, getInputPeer(m.Obj), limit, offset...)
}

type TelegramFileInfo struct {
	InputLocation mtproto.TL
	DcID          int32
	Size          int32
	FName         string
}

// findBestPhotoSize returns largest photo size of images.
// Usually it is the last size-object. But SOME TIMES Sizes aray is reversed.
func findBestPhotoSize(photo mtproto.TL_photo) *mtproto.TL_photoSize {
	var bestSize *mtproto.TL_photoSize
	for _, sizeTL := range photo.Sizes {
		if size, ok := sizeTL.(mtproto.TL_photoSize); ok {
			if bestSize == nil || size.Size > bestSize.Size {
				bestSize = &size
			}
		}
	}
	return bestSize
}

func tgGetMessageMediaFileInfo(msgTL mtproto.TL) *TelegramFileInfo {
	msg, ok := msgTL.(mtproto.TL_message)
	if !ok {
		return nil
	}
	switch media := msg.Media.(type) {
	case mtproto.TL_messageMediaPhoto:
		if _, ok := media.Photo.(mtproto.TL_photoEmpty); ok {
			// Lg.Error("got 'photoEmpty' in media of message #", msg.ID)
			return nil
		}
		photo := media.Photo.(mtproto.TL_photo)
		size := findBestPhotoSize(photo)
		if size == nil {
			// Lg.Error("could not found suitable image size of message #", msg.ID)
			Panicerr("image size search failed")
		}
		return &TelegramFileInfo{
			InputLocation: mtproto.TL_inputPhotoFileLocation{
				ID:            photo.ID,
				AccessHash:    photo.AccessHash,
				FileReference: photo.FileReference,
				ThumbSize:     size.Type,
			},
			Size:  size.Size,
			DcID:  photo.DcID,
			FName: "photo.jpg",
		}
	case mtproto.TL_messageMediaDocument:
		doc := media.Document.(mtproto.TL_document)
		fname := ""
		for _, attrTL := range doc.Attributes {
			if nameAttr, ok := attrTL.(mtproto.TL_documentAttributeFilename); ok {
				fname = nameAttr.FileName
				break
			}
		}
		return &TelegramFileInfo{
			InputLocation: mtproto.TL_inputDocumentFileLocation{
				ID:            doc.ID,
				AccessHash:    doc.AccessHash,
				FileReference: doc.FileReference,
			},
			Size:  doc.Size,
			DcID:  doc.DcID,
			FName: fname,
		}
	default:
		return nil
	}
}

func historyMessage(tg *tgclient.TGClient, inputPeer mtproto.TL, limit int32, offset ...int32) (resmsgs []*TelegramMessageStruct) {
	params := mtproto.TL_messages_getHistory{
		Peer:  inputPeer,
		Limit: limit, // 当前offset开始往上取的条数
		// OffsetID: 4,     // offset的消息的ID，不设置默认最后。例如如果设置offset的id是4，那就是从id3开始往上获取limit个数据。
	}

	if len(offset) != 0 {
		params.OffsetID = offset[0]
	}

	res := tg.SendSyncRetry(params, time.Second, 0, 30*time.Second)

	// Lg.Debug(res)

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
			AccessHash:  tgUser.AccessHash,
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

		// Lg.Debug("====>MsgMap:", msgMap)

		fileInfo := tgGetMessageMediaFileInfo(msg) // 即使客户端发送一条消息多个图片，这里也会是每个消息一个图片，多个消息就是了

		// 如果一条消息，没有消息内容，没有文件，没有action，就跳过
		if Str(msgMap["Message"]) == "" && fileInfo == nil && !Map(msgMap).Has("Action") {
			// Lg.Debug("====>Msg:", msgMap)
			continue
		}
		tms := &TelegramMessageStruct{
			Message:     Str(msgMap["Message"]),
			Chat:        tgchatdatas[Int64(StringMap(msgMap["PeerID"])["ChatID"])], // ChannelID
			User:        tguserdatas[Int64(StringMap(msgMap["FromID"])["UserID"])],
			Time:        Int64(msgMap["Date"]),
			ReplyMarkup: StringMap(msgMap["ReplyMarkup"]),
			ID:          Int64(msgMap["ID"]),
		}
		if Map(msgMap).Has("Entities") {
			tms.Entities = InterfaceArray(msgMap["Entities"])
		}
		if tms.Chat == nil {
			tms.Chat = tgchatdatas[Int64(StringMap(msgMap["PeerID"])["ChannelID"])]
		}
		if fileInfo != nil {
			tms.File = fileInfo
		}
		if Map(msgMap).Has("Action") {
			tms.Action = String(StringMap(msgMap["Action"])["_"]).Replace("TL_messageAction", "").S
		}
		if tms.User == nil {
			if Map(msgMap).Has("PeerID") {
				if Map(StringMap(msgMap["PeerID"])).Has("UserID") {
					tms.User = tguserdatas[Int64(StringMap(msgMap["PeerID"])["UserID"])]
				}
			}
		}
		resmsgs = append(resmsgs, tms)
	}

	return
}

func getInputPeer(Obj mtproto.TL) (inputPeer mtproto.TL) {
	// Lg.Debug(m.Obj)
	switch peer := Obj.(type) {
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
		Peer:     getInputPeer(m.Obj),
		Message:  text,
	}

	res := m.tg.SendSyncRetry(params, time.Second, 0, 30*time.Second)

	// fmt.Println(Repr(res).S)

	if !Repr(res).Has(text) {
		Panicerr("可能发送消息失败？")
	}
}

// --------- get history by username

type TelegramPeerResolved struct {
	tg  *tgclient.TGClient
	Obj mtproto.TL
	// 加上其它peer的信息，例如id，title，username，members之类的
	Type       string // user,channel(group也是channel), 具体是group还是channel需要获取里面的历史消息的时候可以看到
	Name       string // channel和group的话就是title的值，user的话就是FirstName+LastName
	Username   string // 就是ID，可以@的那个
	ID         int64  // 可以用来发起聊天的telegram的这个对象的ID
	AccessHash int64  // 有时候手动构造inputpeer的时候有用, TL_user, TL_channel的时候， TL_chat不用这个。
}

func (m *TelegramStruct) ResolvePeerByUsername(username string) *TelegramPeerResolved {
	params := mtproto.TL_contacts_resolveUsername{
		Username: username,
	}

	r := m.tg.SendSyncRetry(params, time.Second, 0, 30*time.Second)

	// Lg.Debug(r)

	rs := r.(mtproto.TL_contacts_resolvedPeer)

	tgpr := &TelegramPeerResolved{
		tg: m.tg,
	}
	if len(rs.Chats) != 0 {
		tgpr.Obj = rs.Chats[0]
		tgpr.Type = "channel"
		// 有测试group和channel都是TL_channel
		rss := rs.Chats[0].(mtproto.TL_channel)
		tgpr.AccessHash = rss.AccessHash
		tgpr.Name = rss.Title
		tgpr.Username = username
		tgpr.ID = rss.ID
	} else if len(rs.Users) != 0 {
		tgpr.Obj = rs.Users[0]
		tgpr.Type = "user"

		rss := rs.Users[0].(mtproto.TL_user)
		tgpr.AccessHash = rss.AccessHash
		tgpr.Name = rss.FirstName + " " + rss.LastName
		tgpr.Username = username
		tgpr.ID = rss.ID
	} else {
		Panicerr("未找到这个username相关的peer")
		return nil
	}
	return tgpr
}

// 不设置offset则获取最新的limit条数据
// 设置了offset的话，例如如果设置offset的id是4，那就是从id小于4的message开始往上获取limit个数据。（不包括id为4的message）
func (m *TelegramPeerResolved) History(limit int32, offset ...int32) (resmsgs []*TelegramMessageStruct) {
	return historyMessage(
		m.tg,
		getInputPeer(
			m.Obj,
		),
		limit,
		offset...,
	)
}
