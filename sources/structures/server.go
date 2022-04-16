package structures

import (
	"net/http/httputil"
	"sync"
)

type Server struct {
	URL string `json:"url"`
	//URL          *url.URL
	IsAlive      bool
	Mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

func (s *Server) GetAliveStatus() (alive bool) {
	s.Mux.RLock()
	alive = s.IsAlive
	s.Mux.RUnlock()
	return
}

func (s *Server) SetAliveStatus(alive bool) {
	s.Mux.Lock()
	s.IsAlive = alive
	s.Mux.Unlock()
}
