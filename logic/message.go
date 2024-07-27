package logic

import (
	"fmt"
	"time"
)

const (
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

	//Users map[string]*User `json:"users"`
}

func NewMessage(user *User, content string) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeNormal,
		Content: content,
		MsgTime: time.Now(),
	}
}

func NewWelcomeMessage(user *User) *Message {
	return &Message{
		User:    system,
		Type:    MsgTypeSystem,
		Content: fmt.Sprintf("%s 你好，欢迎加入聊天室。", user.NickName),
		MsgTime: time.Now(),
	}
}

func NewEnterMessage(user *User) *Message {
	return &Message{
		User:    system,
		Type:    MsgTypeSystem,
		Content: fmt.Sprintf("%s 加入了聊天室", user.NickName),
		MsgTime: time.Now(),
	}
}

func NewLeaveMessage(user *User) *Message {
	return &Message{
		User:    system,
		Type:    MsgTypeSystem,
		Content: fmt.Sprintf("%s 离开了聊天室", user.NickName),
		MsgTime: time.Now(),
	}
}

func NewErrorMessage(content string) *Message {
	return &Message{
		User:    system,
		Type:    MsgTypeError,
		Content: content,
		MsgTime: time.Now(),
	}
}
