package connect

import (
	"encoding/json"
	"gochat/tools"
)

type Notification struct {
	Type    string      `json:"type"`
	Content string      `json:"content"`
	Data    interface{} `json:"data,omitempty"`
}

func SendNotification(userID string, notification Notification) error {
	msgBytes, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	return tools.PublishNotification(userID, string(msgBytes))
}

// StartNotificationListener 启动通知监听器
func StartNotificationListener(userID string, ch *Channel) {
	pubsub := tools.SubscribeNotification(userID)
	defer pubsub.Close()

	msgChan := pubsub.Channel()
	for msg := range msgChan {
		ch.Send([]byte(msg.Payload))
	}
}
