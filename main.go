package main

import (
	"chat/server"
	"chat/service"
	"log"
)

func main() {
	service.NewService().Run()
	if err := server.ListenAndServe(""); err != nil {
		log.Fatal(err)
	}
}
