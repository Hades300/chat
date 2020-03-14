package handler

import (
	"chat/common"
	"chat/server/api/rpc"
	"github.com/gin-gonic/gin"
)

var userServiceClient = rpc.GetClient()

func Login(c *gin.Context) {
	userForm := common.LoginForm{}
	c.BindJSON(&userForm)
	code, token := userServiceClient.Login(userForm)
	c.JSON(200, common.LoginResult{
		Code:      code,
		AuthToken: token,
	})
}

func Register(c *gin.Context) {
	userForm := common.RegisterForm{}
	c.BindJSON(&userForm)
	code, message := userServiceClient.Register(userForm)
	c.JSON(200, common.RegisterResult{
		Code:    code,
		Message: message,
	})
}
