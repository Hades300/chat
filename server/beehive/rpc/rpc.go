package rpc

import (
	"chat/common"
	"context"
	"github.com/smallnest/rpcx/client"
	"sync"
)

type UserServiceRpcClient struct{}

var Once sync.Once
var userServiceClient = UserServiceRpcClient{}
var d = client.NewPeer2PeerDiscovery("tcp@"+common.UserSericeAddress, "")
var xclient = client.NewXClient("UserServiceRpcServer", client.Failtry, client.RandomSelect, d, client.DefaultOption)

/** 初始化服务发现配置 初始化客户端**/

// 总是返回这个初始化好的
func GetClient() UserServiceRpcClient {
	return userServiceClient
}

func (u *UserServiceRpcClient) CheckAuthByToken(token string) string {
	request := common.CheckAuthRequest{Token: token}
	reply := common.CheckAuthReply{}
	xclient.Call(context.Background(), "CheckAuth", request, &reply)
	return reply.UserName
}
