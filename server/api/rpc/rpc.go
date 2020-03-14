package rpc

import (
	"chat/common"
	"context"
	"github.com/smallnest/rpcx/client"
	"log"
	"sync"
)

type UserServiceRpcClient struct{}

var xclient client.XClient
var Once sync.Once
var userServiceClient = UserServiceRpcClient{}

/** 初始化服务发现配置 初始化客户端**/
// TODO：待会放到api.go::run中初始化
func InitUserSericeRpcClient() {
	d := client.NewPeer2PeerDiscovery("tcp@"+common.UserSericeAddress, "")
	Once.Do(func() {
		xclient = client.NewXClient("UserServiceRpcServer", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	})
}

// 总是返回这个初始化好的
func GetClient() UserServiceRpcClient {
	return userServiceClient
}

func (u *UserServiceRpcClient) Login(userForm common.LoginForm) (Code int, Token string) {
	request := common.LoginRequest{
		UserName: userForm.UserName,
		PassWord: userForm.PassWord,
	}
	reply := common.LoginReply{}
	xclient.Call(context.Background(), "Login", request, &reply)
	Code = reply.Code
	Token = reply.AuthToken
	return
}

func (u *UserServiceRpcClient) Register(userForm common.RegisterForm) (Code int, Message string) {
	request := common.RegisterRequest{
		UserName: userForm.UserName,
		PassWord: userForm.PassWord,
	}
	reply := common.RegisterReply{}
	if err := xclient.Call(context.Background(), "Register", request, &reply); err != nil {
		log.Println("err rpc call", err)
	}
	Code = reply.Code
	Message = reply.Message
	return
}
