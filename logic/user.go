package logic

import (
	"context"
	"errors"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

type User struct {
	conn        *websocket.Conn
	UID         int           `json:"uid"`
	NickName    string        `json:"nickname"`
	EnterAt     time.Time     `json:"enter_at"`
	Addr        string        `json:"addr"`
	MessageChan chan *Message `json:"-"`
}

func NewUser(conn *websocket.Conn) *User {
	return &User{
		conn:        conn,
		MessageChan: make(chan *Message),
	}
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
		_ = NewMessage(u, receiveMsg["content"])
	}
}
