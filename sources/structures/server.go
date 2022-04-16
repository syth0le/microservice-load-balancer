package structures

import (
	"net/http/httputil"
	"sync"
)

type Server struct {
	URL string `json:"url"`
	//URL          *url.URL
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
