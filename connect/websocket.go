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
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	ch := NewChannel(conn, token)
	AddChannel(ch)
	ch.ReadPump()

}
