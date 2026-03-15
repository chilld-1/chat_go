package connect

import (
	"context"
	"gochat/proto"
	"gochat/tools"
	"net"

	"google.golang.org/grpc"
)

type ConnectService struct {
	proto.UnimplementedConnectServiceServer
}

func (s *ConnectService) Broadcast(ctx context.Context, req *proto.BroadcastRequest) (*proto.BroadcastResponse, error) {
	BroadcastMessage([]byte(req.Message))
	return &proto.BroadcastResponse{}, nil
}

func (s *ConnectService) PushPrivateMessage(ctx context.Context, req *proto.PushPrivateMessageRequest) (*proto.PushPrivateMessageResponse, error) {
	SendPrivateMessage(req.UserId, []byte(req.Message))
	return &proto.PushPrivateMessageResponse{}, nil
}
func StartGrpcServer() error {
	// 创建 gRPC 服务器
	s := grpc.NewServer()

	// 注册服务
	proto.RegisterConnectServiceServer(s, &ConnectService{})

	// 启动服务发现
	go tools.RegisterServer("ConnectService", "127.0.0.1:6901")

	// 监听端口
	addr := "127.0.0.1:6901"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// 启动服务器
	return s.Serve(lis)
}
