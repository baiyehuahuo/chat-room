package logic

import (
	"chatroom/global"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var globalUID uint32 = 0
var invalidToken = errors.New("invalid token")
var system = &User{}

type User struct {
	conn        *websocket.Conn
	UID         uint32        `json:"uid"`
	NickName    string        `json:"nickname"`
	EnterAt     time.Time     `json:"enter_at"`
	Addr        string        `json:"addr"`
	MessageChan chan *Message `json:"-"`

	Token string `json:"token"`
	IsNew bool   `json:"is_new"`
}

func NewUser(conn *websocket.Conn, addr, nickname, token string) *User {
	user := &User{
		conn:        conn,
		NickName:    nickname,
		EnterAt:     time.Now(),
		Addr:        addr,
		MessageChan: make(chan *Message),
		Token:       token,
	}

	if user.Token != "" {
		uid, err := parseTokenAndValidate(token, nickname)
		if err == nil {
			user.UID = uid
		}
	}

	if user.UID == 0 {
		user.UID = atomic.AddUint32(&globalUID, 1)
		user.Token = genToken(user.UID, nickname)
		user.IsNew = true
	}

	return user
}

// token 编码
func genToken(uid uint32, nickname string) string {
	message := fmt.Sprintf("%s%s%d", nickname, global.TokenSecret, uid)
	messageMAC := macSha256([]byte(message), []byte(global.TokenSecret))

	return fmt.Sprintf("%suid%d", base64.StdEncoding.EncodeToString(messageMAC), uid)
}

func macSha256(message, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write(message)
	return mac.Sum(nil)
}

// token 解码
func parseTokenAndValidate(token string, nickname string) (uint32, error) {
	pos := strings.LastIndex(token, "uid")
	if pos == -1 {
		return 0, invalidToken
	}

	messageMAC, err := base64.StdEncoding.DecodeString(token[:pos])
	if err != nil {
		return 0, invalidToken
	}

	uid, err := strconv.Atoi(token[pos+len("uid"):])
	if err != nil {
		return 0, invalidToken
	}

	message := fmt.Sprintf("%s%s%d", nickname, global.TokenSecret, uid)
	if message != string(messageMAC) {
		return 0, invalidToken
	}

	if !validateMAC([]byte(message), messageMAC, []byte(global.TokenSecret)) {
		return 0, invalidToken
	}

	return uint32(uid), nil
}

func validateMAC(message, messageMAC, secret []byte) bool {
	mac := hmac.New(sha256.New, secret)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func (u *User) SendMessage(ctx context.Context) {
	for msg := range u.MessageChan {
		err := wsjson.Write(ctx, u.conn, msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func (u *User) ReceiveMessage(ctx context.Context) error {
	var (
		receiveMsg map[string]string
		err        error
	)

	for {
		err = wsjson.Read(ctx, u.conn, &receiveMsg)
		if err != nil {
			if errors.As(err, &websocket.CloseError{}) {
				return nil
			}
			return err
		}
		msg := NewMessage(u, receiveMsg["content"])
		reg := regexp.MustCompile(`@[^\s@]{2,20}`)
		msg.Ats = reg.FindAllString(msg.Content, -1)
		for i := range msg.Ats {
			msg.Ats[i] = msg.Ats[i][1:]
		}
		msg.Content = FilterSensitive(msg.Content)
		Broadcaster.Broadcast(msg)
	}
}
