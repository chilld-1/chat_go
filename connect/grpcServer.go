package connect

import (
	"log"
	"net"

	"github.com/gin-gonic/gin"
)

type GrpcConnect struct {
}

func (c *Connect) RungrpcWebsocket() {
	go StartGrpcServer()
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	r.GET("/ws", WebSocketHandler)

	r.Run(":7000")
}

func (c *Connect) RungrpcTcp() {
	go StartGrpcServer()
	listener, err := net.Listen("tcp", ":7001")
	if err != nil {
		log.Printf("TCP 监听失败: %v", err)
		return
	}
	defer listener.Close()
	log.Println("TCP 服务启动在 :7001")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("TCP 接受连接失败: %v", err)
			continue
		}
		go handleTcpConn(conn)
	}
}
