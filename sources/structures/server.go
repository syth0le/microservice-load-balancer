package structures

import (
	"context"
	"github.com/syth0le/microservice-load-balancer/config"
	"github.com/syth0le/microservice-load-balancer/sources/balancer"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type Server struct {
	URL          string `json:"url"`
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

func (s *Server) AddReverseProxy() {
	serverUrl, err := url.Parse(s.URL)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(s.URL)
	s.ReverseProxy = httputil.NewSingleHostReverseProxy(serverUrl)
	s.ReverseProxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
		log.Printf("[%s] %s\n", serverUrl.Host, e.Error())
		retries := balancer.GetRetryFromContext(request)
		if retries < 3 {
			select {
			case <-time.After(10 * time.Millisecond):
				ctx := context.WithValue(request.Context(), balancer.Retry, retries+1)
				s.ReverseProxy.ServeHTTP(writer, request.WithContext(ctx))
			}
			return
		}

		// after 3 retries, mark this backend as down
		config.ServerPool.MarkBackendStatus(serverUrl, false)

		// if the same request routing for few attempts with different backends, increase the count
		attempts := balancer.GetAttemptsFromContext(request)
		log.Printf("%s(%s) Attempting retry %d\n", request.RemoteAddr, request.URL.Path, attempts)
		ctx := context.WithValue(request.Context(), balancer.Attempts, attempts+1)
		balancer.LoadBalancing(writer, request.WithContext(ctx))
	}
}

func CheckServerAvailability(uri string) bool {
	timeout := 2 * time.Second
	serverUrl, err := url.Parse(uri)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialTimeout("tcp", serverUrl.Host, timeout)
	if err != nil {
		log.Println("Site unreachable, error: ", err)
		return false
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("Cannot close connection %s", serverUrl.Host)
		}
	}(conn)
	return true
}
