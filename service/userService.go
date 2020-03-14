package service

import (
	"chat/common"
	"github.com/smallnest/rpcx/server"
	"chat/service/dao"
	"chat/tool"
	"context"
)

type UserServiceRpcServer struct{}

var sService = sessService{}
var uDao = dao.UserDao{}

/** 业务逻辑： 先查询用户名是否存在，不存在则存入userdata,**/
func (uService *UserServiceRpcServer) Login(ctx context.Context, request common.LoginRequest, reply *common.LoginReply) error {
	user := dao.User{
		UserName: request.UserName,
		PassWord: request.PassWord,
	}
	if uDao.ValidateUserNameAndPassWord(user) {
		token := tool.GetRandomToken(32)
		sService.RedisSaveUser(token, user)
		reply.Code = 200
		reply.AuthToken = token
	} else {
		reply.Code = 400
	}
	return nil
}

func (uService *UserServiceRpcServer) Register(ctx context.Context, request common.RegisterRequest, reply *common.RegisterReply) error {
	user := dao.User{
		UserName: request.UserName,
		PassWord: request.PassWord,
	}
	if uDao.CheckHaveUserName(user.UserName) {
		reply.Code = 400
		reply.Message = "UserName already Used"
	} else {
		uDao.AddUser(user)
		reply.Code = 200
		reply.Message = "注册成功"
	}
	return nil
}

func (uService *UserServiceRpcServer) CheckAuth(ctx context.Context, request common.CheckAuthRequest, reply *common.CheckAuthReply) error {
	user := sService.RedisCheckUserByToken(request.Token)
	if user.UserName == "" {
		reply.Status = 400
		reply.UserName = ""
	} else {
		reply.Status = 200
		reply.UserName = user.UserName
	}
	return nil
}

func initUserServiceRpcServer() {
	s := server.NewServer()
	s.Register(new(UserServiceRpcServer), "")
	go s.Serve("tcp", common.UserSericeAddress)
}
