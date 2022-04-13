package main

import (
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
)

type Server struct {
	URL          *url.URL
	IsAlive      bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

func (s *Server) GetAliveStatus() (alive bool) {
	s.mux.RLock()
	alive = s.IsAlive
	s.mux.RUnlock()
	return
}

func (s *Server) SetAliveStatus(alive bool) {
	s.mux.Lock()
	s.IsAlive = alive
	s.mux.Unlock()
}

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

func main() {

}
