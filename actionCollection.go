package Onyx

import (
	"errors"
	"fmt"
	pb "github.com/go-gems/Onyx/stream"
	"sync"
)

type IncomingMessage struct {
	*pb.Instruction
}

type actionCollection struct {
	collection map[string]*func(stream *Stream, in *IncomingMessage) error
	mux        *sync.Mutex
}

func (c *actionCollection) On(action string, callback func(slave *Stream, in *IncomingMessage) error) {
	c.mux.Lock()
	c.collection[action] = &callback
	c.mux.Unlock()

}

func (c *actionCollection) Off(action string) {
	c.mux.Lock()
	delete(c.collection, action)
	c.mux.Unlock()

}

func (c *actionCollection) do(slave *Stream, instruction *pb.Instruction) error {
	debugPrintf("%v instruction from %v", instruction.Action, instruction.From)

	c.mux.Lock()
	if _, ok := c.collection[instruction.Action]; !ok {
		return errors.New(fmt.Sprintf("action %v not handled", instruction.Action))
	}
	action := *c.collection[instruction.Action]
	c.mux.Unlock()
	return action(slave, &IncomingMessage{Instruction: instruction})
}
func newActionCollection() *actionCollection {
	return &actionCollection{collection: map[string]*func(stream *Stream, in *IncomingMessage) error{}, mux: &sync.Mutex{}}
}
