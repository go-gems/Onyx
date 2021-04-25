package Onyx

import (
	pb "github.com/go-gems/Onyx/stream"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"io"
	"math/rand"
	"net"
	"time"
)

type master struct {
	*actionCollection
	*connectionCollection
	Id       string
	Address  string
	Protocol string
}

func (masterNode *master) Connect(server pb.StreamService_ConnectServer) error {
	debugPrintf("connection established : %v",masterNode.Id)
	server.Send(&pb.Instruction{
		From:  masterNode.Id,
		Action:    "_",
	})
	for {
		instruction, err := server.Recv()
		if err == io.EOF || instruction == nil {
			debugPrintln("empty instruction")
			return nil
		}
		debugPrintf("instruction received : %v - %v", instruction.From, instruction.Action)
		slave := masterNode.fetchOrCreate(instruction.From, masterNode.Id, server, nil)

		if err != nil {
			masterNode.remove(instruction.From)
			debugPrintln(err)
			return nil
		}
		if err = masterNode.do(slave, instruction); err != nil {
			return err
		}
	}

}

func NewMaster(address string) *master {
	rand.Seed(time.Now().UnixNano())

	return &master{
		Id:                   uuid.New().String(),
		Address:              address,
		Protocol:             "tcp",
		connectionCollection: newConnectionCollection(),
		actionCollection:     newActionCollection(),
	}
}
func NewDefaultMaster() *master {
	return NewMaster(":50000")
}

func (masterNode *master) Serve() error {
	lis, err := net.Listen(masterNode.Protocol, masterNode.Address)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	pb.RegisterStreamServiceServer(s, masterNode)
	debugPrintln("starting server")
	return s.Serve(lis)
}
