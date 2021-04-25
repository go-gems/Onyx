package Onyx

import (
	"context"
	pb "github.com/go-gems/Onyx/stream"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"io"
	"log"
)

type Slave struct {
	Id string
	*connectionCollection
	*actionCollection
}

func NewSlave() *Slave {
	s := &Slave{
		Id:                   uuid.New().String(),
		connectionCollection: newConnectionCollection(),
		actionCollection:     newActionCollection(),
	}
	s.On("_", func(slave *Stream, in *IncomingMessage) error {
		return nil
	})

	return s

}

func (s Slave) Join(masterAddress string) error {
	conn, err := grpc.Dial(masterAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}
	client := pb.NewStreamServiceClient(conn)
	stream, err := client.Connect(context.Background())
	if err != nil {
		return err
	}
	go func() {
		for {
			instruction, err := stream.Recv()
			if err == io.EOF || instruction == nil {
				debugPrintln("empty instruction")
			log.Println(err)
			}
			master := s.fetchOrCreate(instruction.From, s.Id, nil, stream)
			if err != nil {
				log.Fatalln(err)
				return
			}
			if err = s.do(master, instruction)
				err != nil {
				log.Fatalln(err)
			}
		}
	}()
	return nil
}

