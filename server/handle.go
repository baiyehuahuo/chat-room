package server

import (
	"chatroom/logic"
	"net/http"
)

func RegisterHandle() {
	go logic.Broadcaster.Start()

	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/ws", WebSocketHandleFunc)
}
