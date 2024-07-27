package logic

import "sync"

type broadcaster struct {
	userLock sync.Mutex
	users    map[string]*User

	enterChannel   chan *User // 用户进入聊天室
	leaveChannel   chan *User // 用户退出聊天室
	messageChannel chan *User // 用户发送普通消息
}

func (b *broadcaster) CanEnterRoom(nickname string) bool {
	b.userLock.Lock()
	defer b.userLock.Unlock()
	_, ok := b.users[nickname]
	return !ok
}
