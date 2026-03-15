package handler

import (
	"gochat/api/grpc"
	"gochat/tools"

	"github.com/gin-gonic/gin"
)

type PushRequest struct {
	Message string `json:"message" binding:"required"`
}

func PushMessage(c *gin.Context) {
	var req PushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequestResponse(c, "参数错误")
		return
	}

	// 调用 gRPC 推送消息
	reply, err := grpc.PushMessage(req.Message)
	if err != nil {
		tools.InternalServerErrorResponse(c, "消息发送失败")
		return
	}

	if reply.Error != "" {
		tools.InternalServerErrorResponse(c, reply.Error)
		return
	}

	tools.SuccessResponse(c, nil)
}
