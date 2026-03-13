package connect

import (
	"encoding/json"
	"gochat/model"
	"gochat/pkg/stickpackage"
	"gochat/tools"
	"log"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// 写入超时
	writeWait = 10 * time.Second

	// 读取超时
	pongWait = 60 * time.Second

	// 发送 ping 的间隔，必须小于 pongWait
	pingPeriod = (pongWait * 9) / 10

	// 最大消息大小
	maxMessageSize = 512
)

type Channel struct {
	conn  *websocket.Conn
	token string
	send  chan []byte
}
type TcpChannel struct {
	// TCP 连接
	conn *net.Conn

	// 用户 token
	token string

	// 消息发送通道
	send chan []byte
}
type Conn interface {
	ReadPump()
	WritePump()
	handleMessage()
	Send()
}

func NewChannel(conn *websocket.Conn, token string) *Channel {
	return &Channel{
		conn:  conn,
		token: token,
		send:  make(chan []byte, 256),
	}
}
func (ch *Channel) ReadPump() {
	defer func() {
		RemoveChannel(ch.token)
		ch.conn.Close()
	}()
	ch.conn.SetReadLimit(maxMessageSize)
	ch.conn.SetReadDeadline(time.Now().Add(pongWait))
	ch.conn.SetPongHandler(func(string) error {
		ch.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := ch.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				tools.Log(err.Error())
			}
			break
		}

		// 解析消息为 model.Message 结构体
		var msg model.Message
		if err := json.Unmarshal(message, &msg); err != nil {
			// 消息格式错误，返回错误信息
			errorMsg := model.Message{
				Type:      "error",
				From:      "system",
				To:        "",
				Content:   "消息格式错误",
				Timestamp: time.Now().Unix(),
			}
			errorBytes, _ := json.Marshal(errorMsg)
			ch.send <- errorBytes
			continue
		}

		// 确保消息字段完整
		if msg.Timestamp == 0 {
			msg.Timestamp = time.Now().Unix()
		}

		ch.handleMessage(msg)
	}
}
func (ch *Channel) handleMessage(message model.Message) {
	// 确保消息字段完整
	if message.Timestamp == 0 {
		message.Timestamp = time.Now().Unix()
	}

	if message.From == "" {
		message.From = ch.token
	}

	// 发送消息到 mq
	msgBytes, err := json.Marshal(message)
	if err != nil {
		// 序列化失败，返回错误
		ch.send <- []byte(`{"type":"error","content":"消息处理失败"}`)
		return
	}

	if err := tools.SendMessage(msgBytes); err != nil {
		ch.send <- []byte(`{"type":"error","content":"消息发送失败"}`)
		return
	}
}

func (ch *Channel) Send(message []byte) {
	select {
	case ch.send <- message:
		// 消息发送成功
	default:
		// 通道已满，关闭连接
		log.Printf("通道已满，关闭连接: %s", ch.token)
		close(ch.send)
		RemoveChannel(ch.token)
	}
}
func (ch *Channel) writerPump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		ch.conn.Close()
	}()

	for {
		select {
		case message, ok := <-ch.send:
			ch.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// 通道已关闭
				ch.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := ch.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 发送队列中的所有消息
			n := len(ch.send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-ch.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			ch.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ch.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func NewTcpChannel(conn net.Conn, token string) *TcpChannel {
	return &TcpChannel{
		conn:  &conn,
		token: token,
		send:  make(chan []byte, 256),
	}
}
func (ch *TcpChannel) ReadPump() {
	defer func() {
		RemoveChannel(ch.token)
		(*ch.conn).Close()
	}()
	(*ch.conn).SetReadDeadline(time.Now().Add(pongWait))

	for {
		// 读取消息长度
		lengthBuf := make([]byte, 4)
		_, err := (*ch.conn).Read(lengthBuf)
		if err != nil {
			if err.Error() != "EOF" {
				log.Printf("TCP 读取错误: %v", err)
			}
			break
		}

		// 解析消息长度
		length := stickpackage.UnpackLength(lengthBuf)

		// 读取消息内容
		dataBuf := make([]byte, length)
		_, err = (*ch.conn).Read(dataBuf)
		if err != nil {
			log.Printf("TCP 读取数据错误: %v", err)
			break
		}
		ch.handleMessageTCP(dataBuf)

		(*ch.conn).SetReadDeadline(time.Now().Add(pongWait))
	}
}

func (ch *TcpChannel) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		(*ch.conn).Close()
	}()
	for {
		select {
		case message, ok := <-ch.send:
			if !ok {
				// 通道已关闭
				return
			}

			// 打包消息
			packedMsg := stickpackage.Pack(message)

			// 发送消息
			_, err := (*ch.conn).Write(packedMsg)
			if err != nil {
				log.Printf("TCP 发送错误: %v", err)
				return
			}
		case <-ticker.C:
			// 发送 ping
			pingMsg := stickpackage.Pack([]byte("ping"))
			_, err := (*ch.conn).Write(pingMsg)
			if err != nil {
				return
			}
		}
	}
}

func (ch *TcpChannel) handleMessageTCP(message []byte) {
	// 解析消息为 model.Message 结构体
	var msg model.Message
	if err := json.Unmarshal(message, &msg); err != nil {
		// 消息格式错误，返回错误信息
		errorMsg := model.Message{
			Type:      "error",
			From:      "system",
			To:        "",
			Content:   "消息格式错误",
			Timestamp: time.Now().Unix(),
		}
		errorBytes, _ := json.Marshal(errorMsg)
		ch.send <- errorBytes
		return
	}

	// 确保消息字段完整
	if msg.Timestamp == 0 {
		msg.Timestamp = time.Now().Unix()
	}

	if msg.From == "" {
		msg.From = ch.token
	}

	// 发送消息到 RabbitMQ
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		ch.send <- []byte(`{"type":"error","content":"消息处理失败"}`)
		return
	}

	if err := tools.SendMessage(msgBytes); err != nil {
		ch.send <- []byte(`{"type":"error","content":"消息发送失败"}`)
		return
	}
}

func (ch *TcpChannel) Send(message []byte) {
	select {
	case ch.send <- message:
		// 消息发送成功
	default:
		// 通道已满，关闭连接
		log.Printf("TCP 通道已满，关闭连接: %s", ch.token)
		close(ch.send)
		RemoveChannel(ch.token)
	}
}
func (ch *TcpChannel) GetToken() string {
	return ch.token
}
