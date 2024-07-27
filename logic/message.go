package logic

import "time"

const (
	_               = iota
	MsgTypeNormal   = iota // 普通消息
	MsgTypeSystem          // 系统消息
	MsgTypeError           // 错误消息
	MsgTypeUserList        // 发送用户列表
)

type Message struct {
	User    *User     `json:"user"`
	Type    int       `json:"type"`
	Content string    `json:"content"`
	MsgTime time.Time `json:"msg_time"`

	Users map[string]*User `json:"users"`
}

func NewMessage(user *User, content string) *Message {}
