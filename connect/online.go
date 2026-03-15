package connect

import (
	"gochat/tools"
	"log"
)

func UpdataUserOnlineStatus(userID string) {
	if err := tools.SetUserOnline(userID); err != nil {
		log.Printf("更新用户在线状态失败: %v", err)
	}
}

func UpdateUserOfflineStatus(userID string) {
	if err := tools.SetUserOffline(userID); err != nil {
		log.Printf("更新用户离线状态失败: %v", err)
	}
}

func CheckUserOnlineStatus(userID string) bool {
	online, err := tools.IsUserOnline(userID)
	if err != nil {
		log.Printf("检查用户在线状态失败: %v", err)
		return false
	}
	return online
}
