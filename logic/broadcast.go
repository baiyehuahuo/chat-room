package logic

import (
	"chatroom/global"
	"log"
	"sync"
)

type broadcaster struct {
	userLock sync.Mutex
	users    map[string]*User

	enterChannel   chan *User    // 用户进入聊天室
	leaveChannel   chan *User    // 用户退出聊天室
	messageChannel chan *Message // 用户发送普通消息
}

var Broadcaster = &broadcaster{
	userLock:       sync.Mutex{},
	users:          make(map[string]*User),
	enterChannel:   make(chan *User),
	leaveChannel:   make(chan *User),
	messageChannel: make(chan *Message, global.MessageQueenLen),
}

func (b *broadcaster) CanEnterRoom(nickname string) bool {
	b.userLock.Lock()
	defer b.userLock.Unlock()
	_, ok := b.users[nickname]
	return !ok
}

func (b *broadcaster) Start() {
	for {
		select {
		case user := <-b.enterChannel:
			b.userLock.Lock()
			b.users[user.NickName] = user
			b.userLock.Unlock()
		case user := <-b.leaveChannel:
			b.userLock.Lock()
			delete(b.users, user.NickName)
			b.userLock.Unlock()
		case message := <-b.messageChannel:
			b.userLock.Lock()
			switch {
			case message.Ats != nil:
				for _, at := range message.Ats {
					if at == message.User.NickName {
						// 不能at自己
						continue
					}

					user, ok := b.users[at]
					if !ok {
						log.Printf("user %s not exist", at)
						continue
					}
					user.MessageChan <- message
				}
			default:
				for _, user := range b.users {
					if user.UID == message.User.UID {
						continue
					}
					user.MessageChan <- message
				}
			}
			b.userLock.Unlock()
		}
	}
}

func (b *broadcaster) UserEntering(u *User) {
	b.enterChannel <- u
}

func (b *broadcaster) UserLeaving(u *User) {
	b.leaveChannel <- u
}

func (b *broadcaster) Broadcast(msg *Message) {
	b.messageChannel <- msg
}
