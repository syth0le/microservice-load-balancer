package sources

import (
	"log"
	"net/url"
	"sync/atomic"
)

type ServerPool struct {
	Servers       []*Server
	CurrentServer uint64
}

func (sp *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&sp.CurrentServer, uint64(1)) % uint64(len(sp.Servers)))
}

func (sp *ServerPool) GetNextServer() *Server {
	next := sp.NextIndex()
	for idx, server := range sp.Servers {
		if server.GetAliveStatus() { // TODO: проверять как то по-другому
			if idx != next {
				atomic.StoreUint64(&sp.CurrentServer, uint64(idx))
			}
			return server
		}
	}
	return nil
}

func (sp *ServerPool) DoHealthCheck() {
	for _, server := range sp.Servers {
		//resp, err := http.Get(server.URL)
		status := CheckServerAvailability(server.URL)
		server.SetAliveStatus(status)

		if !status {
			log.Printf("SERVER CONNECTION REFUSED: %s\n", server.URL)
		} else {
			log.Printf("SERVER: %s is OK\n", server.URL)
		}
	}
}

func (sp *ServerPool) AddServer(server *Server) {
	sp.Servers = append(sp.Servers, server)
}

func (sp *ServerPool) MarkBackendStatus(backendUrl *url.URL, alive bool) {
	for _, b := range sp.Servers {
		if b.URL == backendUrl.String() {
			b.SetAliveStatus(alive)
			break
		}
	}
}
