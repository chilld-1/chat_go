package connect

import "github.com/gin-gonic/gin"

type Connect struct{}

func New() *Connect {
	return &Connect{}
}

func (c *Connect) Run_websocket() {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
	})
	r.GET("/ws", WebSocketHandler)
	r.Run(":7000")
}
