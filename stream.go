package Onyx

import (
	"encoding/json"
	pb "github.com/go-gems/Onyx/stream"
)

type Stream struct {
	remoteId      string
	localId       string
	connectServer pb.StreamService_ConnectServer
	connectClient pb.StreamService_ConnectClient
}

func (s *Stream) GetRemoteId() string {
	return s.remoteId
}
func (s *Stream) GetLocalId() string {
	return s.localId
}

func (s *Stream) Send(action string, item interface{}) error {
	debugPrintf("sending %v from %v to %v", action, s.GetLocalId(), s.GetRemoteId())
	content, err := json.Marshal(item)
	if err != nil {
		return err
	}
	instruction := &pb.Instruction{
		From:    s.localId,
		Action:  action,
		Content: string(content),
	}
	if s.connectServer != nil {
		return s.connectServer.Send(instruction)
	}
	if s.connectClient != nil {
		return s.connectClient.Send(instruction)
	}
	return nil
}
