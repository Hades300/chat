package server

import (
	"chat/server/api/router"
	"net/http"
)

type httpServer struct{}

var DefaultHttpServer *httpServer

func init() {
	DefaultHttpServer = &httpServer{}
}

// singleton pattern && factory

func NewHttpServer() *httpServer {
	return &httpServer{}
}

func (hserver *httpServer) ListenAndServe(addr string) error {
	r := router.Default()
	if addr == "" {
		addr = ":8080"
	}
	return http.ListenAndServe(addr, r)

}
