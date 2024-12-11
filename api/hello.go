1. 1. ```go
      package main
      
      import (
      	"fmt"
      	"net/http"
      	"sync"
      
      	"github.com/gorilla/websocket"
      )
      
      type Server struct {
      	mu      sync.Mutex              // 用于保护 clients 的互斥锁
      	clients map[*websocket.Conn]bool // 已连接的客户端
      	broadcast chan string           // 广播消息的通道
      }
      
      var upgrader = websocket.Upgrader{
      	CheckOrigin: func(r *http.Request) bool {
      		return true // 允许所有跨域连接
      	},
      }
      
      // 启动聊天服务器
      func (s *Server) Start() {
      	s.clients = make(map[*websocket.Conn]bool)
      	s.broadcast = make(chan string)
      
      	// 设置 WebSocket 处理
      	http.HandleFunc("/ws", s.handleConnections)
      
      	// 处理广播消息
      	go s.handleBroadcast()
      
      	// 启动 HTTP 服务
      	port := "8080"
      	fmt.Println("聊天服务器已启动，监听端口:", port)
      	http.ListenAndServe(":"+port, nil)
      }
      
      // 处理 WebSocket 连接
      func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
      	conn, err := upgrader.Upgrade(w, r, nil)
      	if err != nil {
      		fmt.Println("升级到 WebSocket 失败:", err)
      		return
      	}
      	defer conn.Close()
      
      	s.mu.Lock()
      	s.clients[conn] = true
      	s.mu.Unlock()
      
      	fmt.Println("新客户端已连接")
      
      	// 监听客户端消息
      	for {
      		_, message, err := conn.ReadMessage()
      		if err != nil {
      			fmt.Println("客户端断开连接")
      			s.mu.Lock()
      			delete(s.clients, conn)
      			s.mu.Unlock()
      			break
      		}
      		s.broadcast <- string(message)
      	}
      }
      
      // 处理广播，将消息发送给所有客户端
      func (s *Server) handleBroadcast() {
      	for {
      		message := <-s.broadcast
      		s.mu.Lock()
      		for client := range s.clients {
      			err := client.WriteMessage(websocket.TextMessage, []byte(message))
      			if err != nil {
      				fmt.Println("发送消息时出错:", err)
      				client.Close()
      				delete(s.clients, client)
      			}
      		}
      		s.mu.Unlock()
      	}
      }
      
      func main() {
      	server := &Server{}
      	server.Start()
      }
      ```
