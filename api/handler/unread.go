package handler

import (
	"gochat/tools"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func GetUnreadCount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		tools.UnauthorizedResponse(c, "未授权")
		return
	}

	userIDStr := strconv.FormatUint(uint64(userID.(uint)), 10)

	count, err := tools.GetUnreadCount(userIDStr)
	if err != nil {
		if err == redis.Nil {
			// 键不存在，返回 0
			count = 0
		} else {
			// 其他错误，返回 500
			tools.InternalServerErrorResponse(c, "获取未读消息数失败: "+err.Error())
			return
		}
	}

	tools.SuccessResponse(c, gin.H{
		"unread_count": count,
	})
}

func ResetUnreadCount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		tools.UnauthorizedResponse(c, "未授权")
		return
	}

	userIDStr := strconv.FormatUint(uint64(userID.(uint)), 10)

	if err := tools.ResetUnreadCount(userIDStr); err != nil {
		tools.InternalServerErrorResponse(c, "重置未读消息数失败: "+err.Error())
		return
	}

	tools.SuccessResponse(c, nil)
}

func GetRecentMessagesHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		tools.UnauthorizedResponse(c, "未授权")
		return
	}

	userIDStr := strconv.FormatUint(uint64(userID.(uint)), 10)

	countStr := c.DefaultQuery("count", "10")
	count, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		tools.BadRequestResponse(c, "参数错误")
		return
	}

	messages, err := tools.GetRecentMessages(userIDStr, count)
	if err != nil {
		if err == redis.Nil {
			// 键不存在，返回空列表
			messages = []string{}
		} else {
			// 其他错误，返回 500
			tools.InternalServerErrorResponse(c, "获取最近消息失败: "+err.Error())
			return
		}
	}

	tools.SuccessResponse(c, gin.H{
		"messages": messages,
	})
}
