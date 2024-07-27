package logic

import (
	"chatroom/global"
	"container/ring"
)

type offlineProcessor struct {
	n int

	// 保存最近的n条消息
	recentRing *ring.Ring
}

var offlineProcessorInstance = newOfflineProcessor()

func newOfflineProcessor() *offlineProcessor {
	return &offlineProcessor{
		n:          global.OfflineNum,
		recentRing: ring.New(global.OfflineNum),
	}
}

func (o *offlineProcessor) Save(msg *Message) {
	if msg.Type != MsgTypeNormal {
		return
	}
	o.recentRing.Value = msg
	o.recentRing = o.recentRing.Next() // 永远指向没有数据的下一条消息
}

func (o *offlineProcessor) Send(user *User) {
	o.recentRing.Do(func(value any) {
		if value == nil {
			return
		}
		msg := value.(*Message)
		if len(msg.Ats) == 0 { // 发给所有人的
			user.MessageChan <- value.(*Message)
		}

		// 其实这里应该用 uid 才能发给原主，这样应该会被同昵称者接收。
		for _, at := range msg.Ats {
			if at == user.NickName {
				user.MessageChan <- value.(*Message)
			}
		}
	})
}
