package tools

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

const (
	// 用户在线状态前缀
	UserOnlinePrefix = "user:online:"

	// 消息未读计数前缀
	UnreadCountPrefix = "unread:count:"

	// 最近消息缓存前缀
	RecentMessagesPrefix = "recent:messages:"

	// 用户会话前缀
	UserSessionPrefix = "user:session:"

	// 实时通知前缀
	NotificationPrefix = "notification:"

	// 在线状态过期时间
	OnlineExpiration = 30 * time.Minute

	// 会话过期时间
	SessionExpiration = 24 * time.Hour

	// 最近消息缓存数量
	RecentMessagesCount = 100
)

func SetUserOnline(userID string) error {
	key := UserOnlinePrefix + userID
	return RedisClient.Set(key, time.Now().Unix(), OnlineExpiration).Err()
}

// SetUserOffline 设置用户离线状态
func SetUserOffline(userID string) error {
	key := UserOnlinePrefix + userID
	return RedisClient.Del(key).Err()
}

// IsUserOnline 检查用户是否在线
func IsUserOnline(userID string) (bool, error) {
	key := UserOnlinePrefix + userID
	exits, err := RedisClient.Exists(key).Result()
	return exits > 0, err
}

// IncrementUnreadCount 增加消息未读计数
func IncrementUnreadCout(userID string) error {
	key := UnreadCountPrefix + userID
	return RedisClient.Incr(key).Err()
}

// GetUnreadCount 获取消息未读计数
func GetUnreadCount(userID string) (int64, error) {
	key := UnreadCountPrefix + userID
	return RedisClient.Get(key).Int64()
}

// ResetUnreadCount 重置消息未读计数
func ResetUnreadCount(userID string) error {
	key := UnreadCountPrefix + userID
	return RedisClient.Set(key, 0, 0).Err()
}

// AddRecentMessage 添加最近消息
func AddRecentMessage(chatID string, message interface{}) error {
	key := RecentMessagesPrefix + chatID
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// 添加到列表头部
	err = RedisClient.LPush(key, msgBytes).Err()
	if err != nil {
		return err
	}

	// 保持列表长度
	return RedisClient.LTrim(key, 0, RecentMessagesCount-1).Err()
}

// GetRecentMessages 获取最近消息
func GetRecentMessages(chatID string, count int64) ([]string, error) {
	key := RecentMessagesPrefix + chatID
	return RedisClient.LRange(key, 0, count-1).Result()
}

// SetUserSession 设置用户会话
func SetUserSession(userID, sessionData string) error {
	key := UserSessionPrefix + userID
	return RedisClient.Set(key, sessionData, SessionExpiration).Err()
}

// GetUserSession 获取用户会话
func GetUserSession(userID string) (string, error) {
	key := UserSessionPrefix + userID
	return RedisClient.Get(key).Result()
}

// DeleteUserSession 删除用户会话
func DeleteUserSession(userID string) error {
	key := UserSessionPrefix + userID
	return RedisClient.Del(key).Err()
}

// PublishNotification 发布通知
func PublishNotification(userID, notification string) error {
	channel := NotificationPrefix + userID
	return RedisClient.Publish(channel, notification).Err()
}

// SubscribeNotification 订阅通知
func SubscribeNotification(userID string) *redis.PubSub {
	channel := NotificationPrefix + userID
	return RedisClient.Subscribe(channel)
}
