package grpc

import (
	"context"
	"gochat/proto"
	"gochat/tools"
)

func Login(username, password string) (*proto.LoginResponse, error) {
	// 创建 gRPC 客户端
	conn, err := tools.NewGrpcClient("LogicService")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// 创建客户端
	client := proto.NewLogicServiceClient(conn)

	// 调用登录方法
	req := &proto.LoginRequest{
		Username: username,
		Password: password,
	}

	return client.Login(context.Background(), req)
}

// Register 注册
func Register(username, password string) (*proto.RegisterResponse, error) {
	// 创建 gRPC 客户端
	conn, err := tools.NewGrpcClient("LogicService")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// 创建客户端
	client := proto.NewLogicServiceClient(conn)

	// 调用注册方法
	req := &proto.RegisterRequest{
		Username: username,
		Password: password,
	}

	return client.Register(context.Background(), req)
}

// PushMessage 推送消息
func PushMessage(message string) (*proto.PushMessageResponse, error) {
	// 创建 gRPC 客户端
	conn, err := tools.NewGrpcClient("LogicService")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// 创建客户端
	client := proto.NewLogicServiceClient(conn)

	// 调用推送消息方法
	req := &proto.PushMessageRequest{
		Message: message,
	}

	return client.PushMessage(context.Background(), req)
}
