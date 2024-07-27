package logic

import (
	"chatroom/global"
	"container/ring"
)

type offlineProcessor struct {
	n int

	// 保存用户最近的n条消息
	recentRing *ring.Ring

	// 保存某个用户的离线消息
	userRing map[string]*ring.Ring
}

var offlineProcessorInstance = newOfflineProcessor()

func newOfflineProcessor() *offlineProcessor {
	return &offlineProcessor{
		n:          global.OfflineNum,
		recentRing: ring.New(global.OfflineNum),
		userRing:   make(map[string]*ring.Ring),
	}
}

func (o *offlineProcessor) Save(msg *Message) {
	if msg.Type != MsgTypeNormal {
		return
	}
	o.recentRing.Value = msg
	o.recentRing = o.recentRing.Next() // 永远指向没有数据的下一条消息

	for _, nickname := range msg.Ats {
		var r *ring.Ring
		var ok bool
		if r, ok = o.userRing[nickname]; !ok {
			r = ring.New(o.n)
		}
		r.Value = msg
		o.userRing[nickname] = r.Next()
	}
}

func (o *offlineProcessor) Send(user *User) {
	o.recentRing.Do(func(value any) {
		if value != nil {
			user.MessageChan <- value.(*Message)
		}
	})

	if user.IsNew {
		return
	}

	if r, ok := o.userRing[user.NickName]; ok {
		r.Do(func(value any) {
			if value != nil {
				user.MessageChan <- value.(*Message)
			}
		})
		delete(o.userRing, user.NickName)
	}
}
