package connect

import (
	"encoding/json"
	"gochat/model"
	"gochat/tools"
	"log"
)

func StartMessageProcessor() error {
	return tools.ConsumeMessages(func(msgBytes []byte) error {
		var msg model.Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Printf("解析消息失败:%v", err)
			return err
		}
		processMessage(msg)
		return nil
	})
}
func processMessage(msg model.Message) {
	switch msg.Type {
	case "chat":
		handleChatMessage(msg)
	case "notification":
		handleNotificationMessage(msg)
	case "system":
		handleSystemMessage(msg)
	case "broadcast":
		handleBroadcastMessage(msg)
	case "typing":
		handleTypingMessage(msg)
	case "read":
		handleReadMessage(msg)
	default:
		log.Printf("未知消息类型: %s", msg.Type)
	}
}
func handleChatMessage(msg model.Message) {
	if msg.To == "" {
		// 广播消息
		handleBroadcastMessage(msg)
	} else {
		// 私聊消息
		handleSendPrivateMessage(msg)
	}
}

func handleNotificationMessage(msg model.Message) {
	// 查找目标用户并发送通知
	if ch := GetChannel(msg.To); ch != nil {
		msgBytes, _ := json.Marshal(msg)
		switch c := ch.(type) {
		case *Channel:
			c.Send(msgBytes)
		case *TcpChannel:
			c.Send(msgBytes)
		}
	}
}

func handleSystemMessage(msg model.Message) {
	// 系统消息处理逻辑
	tools.Log("系统消息: %s", msg.Content)
	// 可以发送给特定用户或广播
}
func handleSendPrivateMessage(msg model.Message) {
	// 发送私聊消息
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("序列化消息失败: %v", err)
		return
	}
	SendPrivateMessage(msg.To, msgBytes)
}

func handleBroadcastMessage(msg model.Message) {
	// 广播给所有用户
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("序列化消息失败: %v", err)
		return
	}
	BroadcastMessage(msgBytes)
}

func handleTypingMessage(msg model.Message) {
	// 发送输入状态给接收者
	if ch := GetChannel(msg.To); ch != nil {
		msgBytes, _ := json.Marshal(msg)
		switch c := ch.(type) {
		case *Channel:
			c.Send(msgBytes)
		case *TcpChannel:
			c.Send(msgBytes)
		}
	}
}
func handleReadMessage(msg model.Message) {
	// 标记消息为已读
	// 可以更新数据库状态
	log.Printf("消息已读: %s", msg.Content)
}
