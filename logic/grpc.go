package logic

import (
	"context"
	"fmt"
	"gochat/logic/dao"
	"gochat/model"
	"gochat/proto"
	"gochat/tools"
	"net"

	"google.golang.org/grpc"
)

type LogicService struct {
	proto.UnimplementedLogicServiceServer
}

func (s *LogicService) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	// 根据用户名获取用户
	user, err := dao.GetUserByUsername(req.Username)
	if err != nil {
		return &proto.LoginResponse{
			Error: "用户名或密码错误",
		}, nil
	}

	// 验证密码
	if user.Password != req.Password {
		return &proto.LoginResponse{
			Error: "用户名或密码错误",
		}, nil
	}

	// 生成 token
	token, err := tools.GenerateToke(user.ID, user.Username)
	if err != nil {
		return &proto.LoginResponse{
			Error: "生成token失败",
		}, nil
	}

	return &proto.LoginResponse{
		Token:    token,
		Username: user.Username,
	}, nil
}

func (s *LogicService) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	// 检查用户名是否已存在
	_, err := dao.GetUserByUsername(req.Username)
	if err == nil {
		return &proto.RegisterResponse{
			Error: "用户名已存在",
		}, nil
	}

	// 创建用户
	var user = model.User{
		Username: req.Username,
		Password: req.Password,
	}
	err = dao.CreateUser(&user)
	if err != nil {
		return &proto.RegisterResponse{
			Error: "注册失败",
		}, nil
	}

	// 生成 token
	token, err := tools.GenerateToke(user.ID, user.Username)
	if err != nil {
		return &proto.RegisterResponse{
			Error: "生成token失败",
		}, nil
	}

	return &proto.RegisterResponse{
		Token:    token,
		Username: user.Username,
	}, nil
}

// PushMessage 推送消息
func (s *LogicService) PushMessage(ctx context.Context, req *proto.PushMessageRequest) (*proto.PushMessageResponse, error) {
	// 发布消息
	err := tools.SendMessage([]byte(req.Message))
	if err != nil {
		return &proto.PushMessageResponse{
			Error: "消息发送失败",
		}, nil
	}

	return &proto.PushMessageResponse{}, nil
}

func StartGrpcServer() error {
	// 创建 gRPC 服务器
	addr := "127.0.0.1:6900"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s := grpc.NewServer()

	// 注册服务
	proto.RegisterLogicServiceServer(s, &LogicService{})

	// 启动服务发现
	go func() {
		if err := tools.RegisterServer("LogicService", addr); err != nil {
			fmt.Printf("服务发现注册失败: %v\n", err)
		}
	}()

	// 监听端口

	// 启动服务器
	return s.Serve(lis)
}
