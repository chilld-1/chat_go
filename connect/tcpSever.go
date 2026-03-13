package connect

import (
	"encoding/json"
	"errors"
	"gochat/model"
	"gochat/pkg/stickpackage"
	"gochat/tools"
	"log"
	"net"
)

func (c *Connect) RunTcp() {
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

func handleTcpConn(conn net.Conn) {
	defer conn.Close()

	token, err := readToken(conn)
	if err != nil {
		log.Printf("读取 token 失败: %v", err)
		return
	}

	// 验证 token 有效性，与 WebSocket 认证方式保持一致
	if !tools.TokenCheck(nil, token) {
		log.Printf("token 验证失败: %s", token)
		return
	}

	log.Printf("TCP 客户端连接: %s", token)
	tcpCh := NewTcpChannel(conn, token)

	AddChannel(tcpCh)
	tcpCh.ReadPump()

}

func readToken(conn net.Conn) (string, error) {
	lengthBuf := make([]byte, 4)
	_, err := conn.Read(lengthBuf)
	if err != nil {
		log.Printf("读取消息长度失败: %v", err)
		return "", err
	}
	length := stickpackage.UnpackLength(lengthBuf)
	dataBuf := make([]byte, length)
	_, err = conn.Read(dataBuf)
	if err != nil {
		log.Printf("读取消息内容失败: %v", err)
		return "", err
	}

	// 尝试解析 JSON 格式的消息
	var msgTCP model.MessageTCP
	if err := json.Unmarshal(dataBuf, &msgTCP); err == nil && msgTCP.Token != "" {
		// 成功解析 JSON 且包含 token
		log.Printf("从 JSON 消息中提取 token: %s", msgTCP.Token)
		return msgTCP.Token, nil
	}

	// 如果解析失败或没有 token，尝试直接使用消息内容作为 token
	token := string(dataBuf)
	if token == "" {
		return "", errors.New("token 为空")
	}

	log.Printf("直接使用消息内容作为 token: %s", token)
	return token, nil
}
