package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "HTTP, hello")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, req *http.Request) {
		conn, err := websocket.Accept(w, req, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close(websocket.StatusInternalError, "内部出错")

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var v interface{}
		err = wsjson.Read(ctx, conn, &v)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("接收到客户端: %v\n", v)

		err = wsjson.Write(ctx, conn, "Hello Websocket client")
		if err != nil {
			log.Println(err)
			return
		}

		conn.Close(websocket.StatusNormalClosure, "")
	})

	log.Fatal(http.ListenAndServe(":2021", nil))
}
