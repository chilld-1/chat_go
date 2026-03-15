package tools

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RegisterServer(serviceName, addr string) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	defer cli.Close()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		// 设置服务信息
		leaseResp, err := cli.Grant(context.Background(), 30)
		if err != nil {
			return err
		}

		// 注册服务
		serviceKey := fmt.Sprintf("/gochat_srv/%s/%s", serviceName, addr)
		_, err = cli.Put(context.Background(), serviceKey, addr, clientv3.WithLease(leaseResp.ID))
		if err != nil {
			return err
		}

		<-ticker.C
	}
}

func DiscoverService(serviceName string) (string, error) {
	// 创建 etcd 客户端
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return "", err
	}
	defer cli.Close()

	// 查找服务
	servicePrefix := fmt.Sprintf("/gochat_srv/%s/", serviceName)
	resp, err := cli.Get(context.Background(), servicePrefix, clientv3.WithPrefix())
	if err != nil {
		return "", err
	}

	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("service %s not found", serviceName)
	}

	// 返回第一个服务地址
	return string(resp.Kvs[0].Value), nil
}

// NewGrpcClient 创建 gRPC 客户端
func NewGrpcClient(serviceName string) (*grpc.ClientConn, error) {
	// 发现服务
	addr, err := DiscoverService(serviceName)
	if err != nil {
		return nil, err
	}

	// 创建 gRPC 客户端
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
