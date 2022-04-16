package structures

import (
	"log"
	"net/http"
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
		//if true == server.GetAliveStatus() {
		//a = 1
		//}
		resp, err := http.Get(server.URL)

		if err != nil {
			server.SetAliveStatus(false)
			log.Printf(err.Error())
			log.Printf("SERVER CONNECTION REFUSED: %s\n", server.URL)
		} else {
			server.SetAliveStatus(true)
			log.Printf("SERVER: %s is OK - %d\n", server.URL, resp.StatusCode)
		}
	}
}

func (sp *ServerPool) AddServer(server *Server) {
	sp.Servers = append(sp.Servers, server)
}
