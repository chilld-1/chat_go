package router

import (
	"gochat/api/handler"
	"gochat/api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
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
	api := r.Group("/v1")
	{
		user := api.Group("/user")
		{
			user.POST("/login", handler.Login)
			user.POST("/register", handler.Register)
		}
		usergrpc := api.Group("/usergrpc")
		{
			usergrpc.POST("/login", handler.Logingrpc)
			usergrpc.POST("/register", handler.Registergrpc)
		}
		session := api.Group("/session")
		{
			session.Use(middleware.AuthMiddleware())
			session.POST("/set", handler.SetSession)
			session.GET("/get", handler.GetSession)
			session.DELETE("/delete", handler.DeleteSession)
		}

		unread := api.Group("/unread")
		{
			unread.Use(middleware.AuthMiddleware())
			unread.GET("/count", handler.GetUnreadCount)
			unread.POST("/reset", handler.ResetUnreadCount)
		}

		messages := api.Group("/messages")
		{
			messages.Use(middleware.AuthMiddleware())
			messages.GET("/recent", handler.GetRecentMessagesHandler)
		}
	}

	return r
}
