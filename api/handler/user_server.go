package handler

import (
	"gochat/api/grpc"
	"gochat/logic/dao"
	"gochat/model"
	"gochat/tools"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequestResponse(c, err.Error())
		return
	}
	user, err := dao.GetUserByUsername(req.Username)
	if err != nil {
		tools.UnauthorizedResponse(c, "用户名或密码错误")
		return
	}
	if user.Password != req.Password {
		tools.UnauthorizedResponse(c, "用户名或密码错误")
		return
	}
	token, err := tools.GenerateToke(user.ID, user.Username)
	if err != nil {
		tools.InternalServerErrorResponse(c, "生成令牌失败")
		log.Printf("发生错误: %s", err.Error())
		return
	}

	tools.SuccessResponse(c, gin.H{
		"token":    token,
		"username": user.Username,
	})
}

func Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequestResponse(c, "参数错误")
		return
	}
	_, err := dao.GetUserByUsername(req.Username)
	if err == nil {
		// 用户存在
		tools.BadRequestResponse(c, "用户已存在")
		return
	} else if err != gorm.ErrRecordNotFound {
		// 其他错误
		tools.InternalServerErrorResponse(c, "查询用户失败")
		return
	}
	// 用户不存在，可以注册
	err = dao.CreateUser(&model.User{Username: req.Username, Password: req.Password})
	if err != nil {
		tools.InternalServerErrorResponse(c, "注册失败")
		return
	}
	user, err := dao.GetUserByUsername(req.Username)
	token, err := tools.GenerateToke(user.ID, user.Username)
	if err != nil {
		tools.SuccessResponse(c, "已成功创建,请重新登录")
		return
	}
	tools.SuccessResponse(c, gin.H{
		"token":    token,
		"username": user.Username,
	})
}
func Logingrpc(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequestResponse(c, err.Error())
		return
	}
	reply, err := grpc.Login(req.Username, req.Password)
	if err != nil {
		tools.InternalServerErrorResponse(c, "登录失败")
		return
	}

	if reply.Error != "" {
		tools.UnauthorizedResponse(c, reply.Error)
		return
	}

	tools.SuccessResponse(c, gin.H{
		"token":    reply.Token,
		"username": reply.Username,
	})
}
func Registergrpc(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequestResponse(c, "参数错误")
		return
	}

	// 调用 gRPC 注册
	reply, err := grpc.Register(req.Username, req.Password)
	if err != nil {
		tools.InternalServerErrorResponse(c, "注册失败")
		return
	}

	if reply.Error != "" {
		tools.BadRequestResponse(c, reply.Error)
		return
	}

	tools.SuccessResponse(c, gin.H{
		"token":    reply.Token,
		"username": reply.Username,
	})
}
