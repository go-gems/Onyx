package Onyx

import (
	pb "github.com/go-gems/Onyx/stream"
	"sync"
)

type connectionCollection struct {
	collection map[string]*Stream
	mux        *sync.Mutex
}

func newConnectionCollection() *connectionCollection {
	return &connectionCollection{collection: map[string]*Stream{}, mux: &sync.Mutex{}}
}

func (c *connectionCollection) remove(nodeId string) {
	delete(c.collection, nodeId)
}
func (c *connectionCollection) fetch(nodeId string) (*Stream, bool) {
	c.mux.Lock()
	item, ok := c.collection[nodeId]
	c.mux.Unlock()

	return item, ok
}
func (c *connectionCollection) fetchOrCreate(nodeId string, parentNodeId string, server pb.StreamService_ConnectServer, client pb.StreamService_ConnectClient) *Stream {
	c.mux.Lock()
	defer c.mux.Unlock()
	if _, ok := c.collection[nodeId]; !ok {
		c.collection[nodeId] = &Stream{
			remoteId:      nodeId,
			localId:       parentNodeId,
			connectServer: server,
			connectClient: client,
		}
	}
	return c.collection[nodeId]
}

func (c *connectionCollection) Broadcast(action string, content interface{}) error {
	for _, stream := range c.collection {
		if err := stream.Send(action, content); err != nil {
			return err
		}
	}
	return nil
}
