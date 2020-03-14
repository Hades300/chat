package server

import (
	"chat/server/beehive/context"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type wsserver struct{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var beehiveContext = context.NewBeehiveContext()
var task = sync.Once{}

func (ws *wsserver) ListenAndServe(addr string) error {
	http.HandleFunc("/ws", wsHandleFunc)
	return http.ListenAndServe(addr, nil)
}

func wsHandleFunc(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Error occur when upgrade the conn")
	}
	hcon := context.NewDefaultHcon(conn)

	// for debug  test

	task.Do(func() { go beehiveContext.RoomSuperviser() })
	go hcon.Start()
	go hcon.ConsumeMessageQueue()
}
