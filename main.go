package main

import (
	"flag"
	"fmt"
	"gochat/api"
	"gochat/config"
	"gochat/connect"
	"gochat/db"
	"gochat/tools"
)

func main() {
	err := config.Init()
	if err != nil {
		fmt.Printf("初始化配置错误:%v\n", err)
		return
	}
	err = db.Init()
	if err != nil {
		fmt.Printf("数据库初始化错误:%v\n", err)
		return
	}
	defer db.Close()
	err = tools.InitRedis(
		config.AppConfig.CommonRedis.RedisAddress,
		config.AppConfig.CommonRedis.RedisPassword,
		config.AppConfig.CommonRedis.DB,
	)
	if err != nil {
		fmt.Printf("初始化 Redis 失败: %v\n", err)
		return
	}
	defer tools.CloseRedis()

	err = tools.InitRabbitMQ(config.AppConfig.CommonRabbitMQ.URL)
	if err != nil {
		fmt.Printf("初始化 RabbitMQ 失败: %v\n", err)
		return
	}
	defer tools.CloseRabbitMQ()
	if err := connect.StartMessageProcessor(); err != nil {
		fmt.Printf("启动消息处理器失败: %v\n", err)
		return
	}

	var module string
	flag.StringVar(&module, "module", "", "指定运行模块")
	flag.Parse()
	fmt.Printf("启动模块: %s\n", module)
	switch module {
	case "api":
		api.New().Run()
	case "connect_websocket":
		connect.New().Run_websocket()
	case "connect_tcp":
		connect.New().RunTcp()
	default:
		fmt.Println("未知模块")
	}

}
