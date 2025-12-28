package websocket

import (
	"encoding/json"
	"log"

	"go-flutter-mall/backend/config"
	"go-flutter-mall/backend/models"
)

// Hub 维护活跃的客户端连接并将消息广播给客户端
type Hub struct {
	// 已注册的客户端
	Clients map[*Client]bool

	// 注册请求通道
	Register chan *Client

	// 注销请求通道
	Unregister chan *Client

	// 消息广播通道 (这里处理的是业务层面的消息)
	Broadcast chan *models.ChatMessage
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan *models.ChatMessage),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			log.Printf("Client registered: %s (Type: %s, ID: %d)", client.ID, client.Type, client.UserID)

		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				log.Printf("Client unregistered: %s", client.ID)
			}

		case message := <-h.Broadcast:
			// 将消息保存到数据库
			if err := config.DB.Create(message).Error; err != nil {
				log.Printf("Failed to save message: %v", err)
			}

			// 路由消息：找到目标接收者
			log.Printf("Broadcasting message from %s (ID: %d) to ...", message.SenderType, message.SenderID)
			for client := range h.Clients {
				// 逻辑：
				// 1. 如果是用户发给客服 (SenderType=user)，则广播给所有客服 (Type=admin)
				// 2. 如果是客服发给用户 (SenderType=admin)，则发给指定用户 (Type=user, UserID=ReceiverID)
				// 3. 同时也发回给发送者自己 (为了多端同步)

				shouldSend := false

				if client.UserID == message.SenderID && client.Type == message.SenderType {
					// 发给发送者自己
					shouldSend = true
				} else if message.SenderType == "user" {
					// 用户发消息，发给所有客服
					if client.Type == "admin" {
						shouldSend = true
					}
				} else if message.SenderType == "admin" {
					// 客服发消息
					// 1. 发给指定用户
					if client.Type == "user" && client.UserID == message.ReceiverID {
						shouldSend = true
					}
					// 2. 发给其他客服 (同步消息)
					if client.Type == "admin" && client.UserID != message.SenderID {
						shouldSend = true
					}
				}

				if shouldSend {
					log.Printf("Sending to %s (ID: %d)", client.Type, client.UserID)
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(h.Clients, client)
					}
				}
			}
		}
	}
}

// WSMessage 包装 WebSocket 传输的消息结构
type WSMessage struct {
	Type    string              `json:"type"`    // message, heartbeat
	Payload *models.ChatMessage `json:"payload"` // 实际的聊天消息
}

func (h *Hub) HandleMessage(client *Client, msg []byte) {
	var wsMsg WSMessage
	if err := json.Unmarshal(msg, &wsMsg); err != nil {
		log.Printf("Invalid message format: %v", err)
		return
	}

	if wsMsg.Type == "message" && wsMsg.Payload != nil {
		chatMsg := wsMsg.Payload
		// 补全发送者信息
		chatMsg.SenderID = client.UserID
		chatMsg.SenderType = client.Type

		log.Printf("Received message from %s (ID: %d): %s", client.Type, client.UserID, chatMsg.Content)

		h.Broadcast <- chatMsg
	}
}
