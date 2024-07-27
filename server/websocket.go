package server

import (
	"chatroom/logic"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func WebSocketHandleFunc(w http.ResponseWriter, req *http.Request) {
	conn, err := websocket.Accept(w, req, nil)
	if err != nil {
		log.Println("websocket accept error:", err)
		return
	}

	// 检测昵称长度是否合法
	nickname := req.FormValue("nickname")
	if l := len(nickname); l < 4 || l > 20 {
		log.Println("nickname length error", nickname, len(nickname))
		wsjson.Write(req.Context(), conn, logic.NewErrorMessage("非法昵称，昵称长度 4~20"))
		conn.Close(websocket.StatusUnsupportedData, "nickname length error")
		return
	}

	// 校验昵称是否已存在
	if !logic.Broadcaster.CanEnterRoom(nickname) {
		log.Println("nickname exists")
		wsjson.Write(req.Context(), conn, logic.NewErrorMessage("昵称已存在"))
		conn.Close(websocket.StatusUnsupportedData, "nickname exists")
		return
	}

	addr := req.RemoteAddr
	user := logic.NewUser(conn, addr, nickname, "")

	// 启动用户发送消息的线程
	go user.SendMessage(req.Context())

	// 给用户发送欢迎消息
	user.MessageChan <- logic.NewWelcomeMessage(user)

	// 告知当前用户新用户的到来
	logic.Broadcaster.Broadcast(logic.NewEnterMessage(user))

	// 将用户加入用户列表
	logic.Broadcaster.UserEntering(user)
	log.Printf("user: %s, joins chat", user.NickName)

	// 处理用户消息
	err = user.ReceiveMessage(req.Context())

	// 用户离开
	logic.Broadcaster.UserLeaving(user)
	logic.Broadcaster.Broadcast(logic.NewLeaveMessage(user))
	log.Printf("user: %s, leaves chat", user.NickName)

	if err != nil {
		log.Printf("read from client %s error: %v", user.NickName, err)
		conn.Close(websocket.StatusInternalError, "Read from client error")
		return
	}

	conn.Close(websocket.StatusNormalClosure, "")
}
