package main

import (
	"log"

	interaction "github.com/nnieie/golanglab5/kitex_gen/interaction/interactionservice"
)

func main() {
	svr := interaction.NewServer(new(InteractionServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
