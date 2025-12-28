package websocket

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"go-flutter-mall/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client 代表一个 WebSocket 连接
type Client struct {
	Hub *Hub

	// WebSocket 连接
	Conn *websocket.Conn

	// 发送消息的缓冲通道
	Send chan *models.ChatMessage

	// 客户端标识
	ID     string
	UserID uint
	Type   string // "user" or "admin"
}

// readPump 泵送来自 WebSocket 连接的消息到 Hub
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.Hub.HandleMessage(c, message)
	}
}

// writePump 泵送来自 Hub 的消息到 WebSocket 连接
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub 关闭了通道
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// 重新封装
			wsMsg := map[string]interface{}{
				"type":    "message",
				"payload": message,
			}

			if err := c.Conn.WriteJSON(wsMsg); err != nil {
				return
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				msg := <-c.Send
				c.Conn.WriteJSON(map[string]interface{}{
					"type":    "message",
					"payload": msg,
				})
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs 处理 WebSocket 请求
func ServeWs(hub *Hub, c *gin.Context) {
	// 获取用户身份
	// 如果是用户，从 JWT 中获取 (已经在中间件中设置了 userID)
	// 如果是管理员，也从 JWT 中获取

	// 这里假设中间件已经设置了 "userID" (uint) 或 "adminID" (uint)
	// 并且有一个标识 "userType" ("user" or "admin")

	// 为了简化，我们允许 query param 传递 token (因为 WS 连接时不方便带 Header)
	// 或者前端先握手。
	// 这里简单起见，假设 token 放在 query "token" 中，或者由中间件处理

	// 假设中间件已过，直接取
	// 但是 gin 的中间件对于 WS 升级请求可能有点问题，因为 WS 是 GET 请求
	// 通常做法：前端连接 ws://host/ws?token=...
	// 我们手动解析 token

	userID := c.Query("user_id") // 临时方案：直接传 ID (极不安全，仅供演示)
	userType := c.Query("type")  // "user" or "admin"

	// TODO: 验证 token
	// token := c.Query("token")
	// if token == "" { ... }

	// 临时：信任前端传来的 ID
	// 实际项目必须解析 token

	if userID == "" || userType == "" {
		// 尝试从 Context 取 (如果走了中间件)
		// ...
		log.Println("Missing user_id or type")
		// c.Status(http.StatusUnauthorized)
		// WS 握手失败
		return
	}

	uid, _ := strconv.Atoi(userID)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		Hub:    hub,
		Conn:   conn,
		Send:   make(chan *models.ChatMessage, 256),
		ID:     userID, // String ID for logging
		Type:   userType,
		UserID: uint(uid),
	}

	// 注册
	client.Hub.Register <- client

	// 开启协程
	go client.WritePump()
	go client.ReadPump()
}
