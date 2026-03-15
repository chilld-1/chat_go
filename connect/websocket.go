package connect

import (
	"gochat/tools"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketHandler(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		tools.BadRequestResponse(c, "缺失token")
		return
	}
	if !tools.TokenCheck(c, token) {
		return
	}
	userID := c.GetString("user_id")

	// 更新用户在线状态
	UpdataUserOnlineStatus(userID)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	ch := NewChannel(conn, token)
	AddChannel(ch)
	go StartNotificationListener(userID, ch)
	defer UpdateUserOfflineStatus(userID)
	ch.ReadPump()

}
