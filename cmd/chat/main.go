package main

import (
	"log"

	chat "github.com/nnieie/golanglab5/kitex_gen/chat/chatservice"
)

func main() {
	svr := chat.NewServer(new(ChatServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
