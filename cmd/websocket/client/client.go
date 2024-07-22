package main

import (
	"context"
	"fmt"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://localhost:2021/ws", nil)
	if err != nil {
		panic(err)
	}
	defer c.Close(websocket.StatusInternalError, "内部出错")

	err = wsjson.Write(ctx, c, "Hello Websocket server")
	if err != nil {
		panic(err)
	}

	var v interface{}
	err = wsjson.Read(ctx, c, &v)
	if err != nil {
		panic(err)
	}

	fmt.Printf("接收到服务端响应: %v\n", v)
	c.Close(websocket.StatusNormalClosure, "")
}
