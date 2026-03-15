package handler

import (
	"gochat/tools"

	"github.com/gin-gonic/gin"
)

type SessionRequest struct {
	SessionsData string `json:"session_data" binding:"required"`
}

// SetSession 设置会话
func SetSession(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		tools.UnauthorizedResponse(c, "未授权")
		return
	}
	var req SessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequestResponse(c, "参数错误")
		return
	}
	if err := tools.SetUserSession(userID, req.SessionsData); err != nil {
		tools.InternalServerErrorResponse(c, "设置会话失败")
		return
	}
	tools.SuccessResponse(c, nil)
}

// GetSession 获取会话
func GetSession(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		tools.UnauthorizedResponse(c, "未授权")
		return
	}

	sessionData, err := tools.GetUserSession(userID)
	if err != nil {
		tools.InternalServerErrorResponse(c, "获取会话失败")
		return
	}

	tools.SuccessResponse(c, gin.H{"session_data": sessionData})
}

// DeleteSession 删除会话
func DeleteSession(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		tools.UnauthorizedResponse(c, "未授权")
		return
	}

	if err := tools.DeleteUserSession(userID); err != nil {
		tools.InternalServerErrorResponse(c, "删除会话失败")
		return
	}

	tools.SuccessResponse(c, nil)
}
