package utils

import (
	"github.com/syth0le/microservice-load-balancer/config"
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
			config.ServerPool.DoHealthCheck()
			log.Println("Finished Health Check")
		}
	}
}

func CreateServerPool() {
	for _, server := range config.Cfg.Servers {
		server.AddReverseProxy()
		config.ServerPool.AddServer(&structures.Server{
			URL:          server.URL,
			ReverseProxy: server.ReverseProxy,
			IsAlive:      server.IsAlive,
		})
	}
	config.ServerPool.DoHealthCheck()
}
