package server

import (
	"errors"
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
		log.Println("nickname length error")
		wsjson.Write(req.Context(), conn, errors.New("非法昵称，昵称长度 4~20"))
		conn.Close(websocket.StatusUnsupportedData, "nickname length error")
		return
	}

	// todo 校验昵称是否已存在

	// todo 构建用户实例

	// todo 启动用户发送消息的线程

	// todo 给用户发送欢迎消息

	// todo 告知当前用户新用户的到来

	// todo 将用户加入用户列表

	// todo 处理用户消息

	// todo 用户离开

	// todo 处理错误
}
