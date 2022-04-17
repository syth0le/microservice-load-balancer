package utils

import (
	"github.com/syth0le/microservice-load-balancer/sources/structures"
	"log"
	"time"
)

func DoHealthCheck() {
	for {
		countdown := time.NewTicker(time.Minute)
		select {
		case <-countdown.C:
			log.Println("Started Health Check")
			structures.ServPool.DoHealthCheck()
			log.Println("Finished Health Check")
		}
	}
}

func CreateServerPool() {
	for _, server := range structures.Cfg.Servers {
		server.AddReverseProxy()
		structures.ServPool.AddServer(&structures.Server{
			URL:          server.URL,
			ReverseProxy: server.ReverseProxy,
			IsAlive:      server.IsAlive,
		})
	}
	structures.ServPool.DoHealthCheck()
}
