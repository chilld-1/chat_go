package model

type Message struct {
	Type      string `json:"type"`      // 消息类型：chat, notification, system
	From      string `json:"from"`      // 发送者
	To        string `json:"to"`        // 接收者（空表示广播）
	Content   string `json:"content"`   // 消息内容
	Timestamp int64  `json:"timestamp"` // 时间戳
	Data      any    `json:"data"`      // 附加数据
}

type MessageTCP struct {
	Type      string `json:"type"`      // 消息类型：chat, notification, system
	From      string `json:"from"`      // 发送者
	To        string `json:"to"`        // 接收者（空表示广播）
	Content   string `json:"content"`   // 消息内容
	Timestamp int64  `json:"timestamp"` // 时间戳
	Data      any    `json:"data"`      // 附加数据
	Token     string `json:"token"`
}
