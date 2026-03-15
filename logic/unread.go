package logic

import (
	"gochat/tools"
)

func IncrementUnreadCout(userID string) error {
	return tools.IncrementUnreadCout(userID)
}

func GetUnreadCount(userID string) (int64, error) {
	return tools.GetUnreadCount(userID)
}

func ResetUnreadCount(userID string) error {
	return tools.ResetUnreadCount(userID)
}
