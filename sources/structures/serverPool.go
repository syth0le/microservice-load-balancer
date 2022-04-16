package structures

import (
	"log"
	"net/http"
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

func (sp *ServerPool) DoHealthCheck() {
	for _, server := range sp.servers {
		if true == server.GetAliveStatus() {
			//a = 1
		}
		//pingedURL, err := url.Parse(server.URL)
		resp, err := http.Get(server.URL)

		if err != nil {
			//log.Fatal(err.Error())
			log.Printf("SERVER CONNECTION REFUSED: %s - %d\n", server.URL, resp)
		} else {
			server.SetAliveStatus(true)
			log.Printf("SERVER: %s is OK - %d\n", server.URL, resp.StatusCode)
		}
	}
}
