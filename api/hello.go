package handler

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Server struct {
	mu        sync.Mutex              // 用于保护 clients 的互斥锁
	clients   map[*websocket.Conn]bool // 已连接的客户端
	broadcast chan string             // 广播消息的通道
}

var server = Server{
	clients:   make(map[*websocket.Conn]bool),
	broadcast: make(chan string),
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有跨域连接
	},
}

// WebSocket 处理函数，符合 Vercel 的要求
func Handler(w http.ResponseWriter, r *http.Request) {
	// 升级 HTTP 请求为 WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("升级到 WebSocket 失败:", err)
		return
	}
	defer conn.Close()

	// 将连接加入到客户端列表
	server.mu.Lock()
	server.clients[conn] = true
	server.mu.Unlock()

	// 监听客户端消息
	go handleMessages(conn)

	// 广播消息给所有客户端
	go handleBroadcast()
}

func handleMessages(conn *websocket.Conn) {
	defer func() {
		server.mu.Lock()
		delete(server.clients, conn)
		server.mu.Unlock()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("读取客户端消息时出错:", err)
			break
		}
		server.broadcast <- string(message)
	}
}

func handleBroadcast() {
	for {
		message := <-server.broadcast
		server.mu.Lock()
		for client := range server.clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				client.Close()
				delete(server.clients, client)
			}
		}
		server.mu.Unlock()
	}
}
