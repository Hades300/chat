package router

import (
	"chat/server/api/handler"
	"github.com/gin-gonic/gin"
)

func Default() *gin.Engine {
	r := gin.Default()
	initUserRouter(r)
	return r
}

func initUserRouter(r *gin.Engine) {
	userGroup := r.Group("/user")
	{
		userGroup.POST("/login", handler.Login)
		userGroup.POST("/register", handler.Register)
	}
}
