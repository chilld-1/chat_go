package connect

import (
	"sync"
)

var (
	channels = make(map[string]interface{})
	mu       sync.RWMutex
)

func AddChannel(ch interface{}) {
	mu.Lock()
	defer mu.Unlock()

	var token string
	switch c := ch.(type) {
	case *Channel:
		token = c.token
		go c.writerPump()
	case *TcpChannel:
		token = c.token
		go c.WritePump()
	}

	channels[token] = ch

}
func RemoveChannel(token string) {
	mu.Lock()
	defer mu.Unlock()
	if ch, ok := channels[token]; ok {
		switch c := ch.(type) {
		case *Channel:
			close(c.send)
		case *TcpChannel:
			close(c.send)
		}
		delete(channels, token)
	}

}
func GetChannel(token string) interface{} {
	mu.RLock()
	defer mu.RUnlock()
	return channels[token]
}
func BroadcastMessage(message []byte) {
	mu.RLock()
	defer mu.RUnlock()
	for _, ch := range channels {
		switch c := ch.(type) {
		case *Channel:
			c.Send(message)
		case *TcpChannel:
			c.Send(message)
		}
	}
}
func SendPrivateMessage(token string, message []byte) {
	mu.RLock()
	defer mu.RUnlock()
	if ch, ok := channels[token]; ok {
		switch c := ch.(type) {
		case *Channel:
			c.Send(message)
		case *TcpChannel:
			c.Send(message)
		}
	}
}
