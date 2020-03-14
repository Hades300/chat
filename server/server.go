package server

import "chat/server/api"

func ListenAndServe(addr string) error {
	api.Run()
	wsserver := wsserver{}
	hserver := NewHttpServer()
	go wsserver.ListenAndServe(":8081")
	return hserver.ListenAndServe(":8082")
}
