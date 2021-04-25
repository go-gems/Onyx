package main

import (
	"github.com/go-gems/Onyx"
	"log"
	"time"
)

func main() {

	server := Onyx.NewSlave()
	server.On("pong", func(master *Onyx.Stream, in *Onyx.IncomingMessage) error {
		log.Printf("received pong from %v to %v", master.GetRemoteId(), master.GetLocalId())

		log.Println(in.Content)
		return nil
	})

	server.Join(":50001")
	for {
		log.Printf("Broadcasting PING from %v", server.Id)

		server.Broadcast("ping", "Hello there!")

		time.Sleep(100 * time.Millisecond)
	}

}
