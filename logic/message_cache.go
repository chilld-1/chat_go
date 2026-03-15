package logic

import (
	"encoding/json"
	"gochat/model"
	"gochat/tools"
)

func AddRecentMessage(chatID string, message model.Message) error {
	return tools.AddRecentMessage(chatID, message)
}

func GetRecentMessages(chatID string, count int64) ([]model.Message, error) {
	msgStrings, err := tools.GetRecentMessages(chatID, count)
	if err != nil {
		return nil, err
	}
	var message []model.Message
	for _, msgStr := range msgStrings {
		var msg model.Message
		if err := json.Unmarshal([]byte(msgStr), &msg); err != nil {
			continue
		}
		message = append(message, msg)
	}
	return message, nil
}
