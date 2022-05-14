package sources

import (
	"log"
	"time"
)

func DoHealthCheck() {
	for {
		countdown := time.NewTicker(time.Minute)
		select {
		case <-countdown.C:
			log.Println("Started Health Check")
			ServPool.DoHealthCheck()
			log.Println("Finished Health Check")
		}
	}
}

func CreateServerPool() {
	for _, server := range Cfg.Servers {
		server.AddReverseProxy()
		ServPool.AddServer(&Server{
			URL:          server.URL,
			ReverseProxy: server.ReverseProxy,
			IsAlive:      server.IsAlive,
		})
	}
	ServPool.DoHealthCheck()
}
