package main

import (
	"github.com/go-gems/Onyx"
	"log"
)

func main() {

	server := Onyx.NewMaster(":50001")
	server.On("ping", func(slave *Onyx.Stream, in *Onyx.IncomingMessage) error {

		log.Printf("received ping from %v to %v", slave.GetRemoteId(), slave.GetLocalId())
		return slave.Send("pong", "my name is Henri")
	})
	log.Fatalln(server.Serve())
}
