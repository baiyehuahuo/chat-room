package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

var (
	// 新用户到来
	enteringChannel = make(chan *User)
	// 用户离开
	leavingChannel = make(chan *User)
	// 广播用的普通消息 channel
	messageChannel = make(chan Message, 8)
)

func main() {
	listener, err := net.Listen("tcp", ":2020")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConn(conn)
	}

}

// broadcaster 记录聊天室用户并进行广播
// 1. 新用户进来
// 2. 用户普通消息
// 3. 用户离开
func broadcaster() {
	users := make(map[*User]struct{})

	for {
		select {
		case user := <-enteringChannel:
			// 新用户注册
			users[user] = struct{}{}
		case user := <-leavingChannel:
			// 老用户离开
			delete(users, user)
		case msg := <-messageChannel:
			// 用户发送消息
			for user := range users {
				if user.ID == msg.UserID {
					continue
				}
				user.MessageChannel <- msg
			}
		}
	}
}

type User struct {
	ID             int
	Addr           string
	EnterAt        time.Time
	MessageChannel chan Message
}

type Message struct {
	UserID  int
	Content string
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	// 新用户进来 构建实例
	user := &User{
		ID:             rand.Int(), // 假设不会随机到0 是全局信息
		Addr:           conn.RemoteAddr().String(),
		EnterAt:        time.Now(),
		MessageChannel: make(chan Message, 8),
	}
	msg := Message{UserID: user.ID}

	// 新建线程用于读操作读写
	go sendMessage(conn, user.MessageChannel)

	// 发送欢迎消息, 通知所有用户新用户的到来
	user.MessageChannel <- Message{Content: fmt.Sprintf("Welcome, %d!\n", user.ID)}
	messageChannel <- Message{Content: fmt.Sprintf("user: %d has enter.\n", user.ID)}

	// 将用户记录到全局用户列表
	enteringChannel <- user

	// 监测用户活跃情况
	var userActive = make(chan struct{})
	go func() {
		d := time.Minute
		timer := time.NewTimer(d)
		for {
			select {
			case <-timer.C:
				conn.Close() // 没有处理错误，因为出现错误也就是重复关闭而已，已经达到了关闭的目标
			case <-userActive:
				timer.Reset(d)
			}
		}
	}()

	// 读取用户输入
	input := bufio.NewScanner(conn)
	for input.Scan() {
		msg.Content = fmt.Sprintf("%d: %s.\n", user.ID, input.Text())
		messageChannel <- msg
		userActive <- struct{}{}
	}

	if err := input.Err(); err != nil {
		log.Println(err)
	}

	// 用户离开
	leavingChannel <- user
	messageChannel <- Message{Content: fmt.Sprintf("user: %d has left.\n", user.ID)}
}

// sendMessage 发送消息
func sendMessage(conn net.Conn, ch <-chan Message) {
	for msg := range ch {
		_, err := fmt.Fprintf(conn, msg.Content)
		if err != nil {
			log.Println(err)
		}
	}
}
