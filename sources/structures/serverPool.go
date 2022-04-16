package structures

import (
	"sync/atomic"
)

type ServerPool struct {
	servers       []*Server
	currentServer uint64
}

func (sp *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&sp.currentServer, uint64(1)) % uint64(len(sp.servers)))
}

func (sp *ServerPool) GetNextServer() *Server {
	next := sp.NextIndex()
	for idx, server := range sp.servers {
		if server.GetAliveStatus() {
			if idx != next {
				atomic.StoreUint64(&sp.currentServer, uint64(idx))
			}
			return server
		}
	}
	return nil
}
