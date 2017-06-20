package gateway

import (
	"netlink"
	"fmt"
	"sync"
)

type ServerEndpoint struct {
	Type 				int
	Id 					int
	Session 			*netlink.Session
}

type Endpoint struct {
	sync.RWMutex
	endpointList 		map[int]*ServerEndpoint
}

func NewEndpoint() *Endpoint {
	return &Endpoint{
		endpointList: make(map[int]*ServerEndpoint),
	}
}

func (ep *Endpoint) Add(sep *ServerEndpoint) error {
	ep.Lock()
	defer ep.Unlock()

	if _, ok := ep.endpointList[sep.Id]; ok {
		return fmt.Errorf("server endpoint id is exists %d", sep.Id)
	}
	ep.endpointList[sep.Id] = sep
	return nil
}

func (ep *Endpoint) Get(id int) *ServerEndpoint {
	ep.RLock()
	defer ep.RUnlock()

	if _, ok := ep.endpointList[id]; !ok {
		return nil
	}
	return ep.endpointList[id]
}

func (ep *Endpoint) Del(id int) error {
	ep.Lock()
	defer ep.Unlock()

	if _, ok := ep.endpointList[id]; ok {
		delete(ep.endpointList, id)
	}

	return  nil
}
