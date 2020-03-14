package api

import "chat/server/api/rpc"

func Run() {
	rpc.InitUserSericeRpcClient()
}
