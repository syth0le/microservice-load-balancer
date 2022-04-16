package main

import (
	"encoding/json"
	"fmt"
	"github.com/syth0le/microservice-load-balancer/config"
	"github.com/syth0le/microservice-load-balancer/sources/structures"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var cfg config.Config
var serverPool structures.ServerPool

func doHealthCheck() {
	for {
		countdown := time.NewTicker(2 * time.Minute)
		select {
		case <-countdown.C:
			log.Println("Started Health Check")
			serverPool.DoHealthCheck()
			log.Println("Finished Health Check")
		}
	}
}
func createServerPool() {
	for _, server := range cfg.Servers {
		server.AddReverseProxy()
		serverPool.AddServer(&structures.Server{
			URL:          server.URL,
			ReverseProxy: server.ReverseProxy,
			IsAlive:      server.IsAlive,
		})
		//log.Printf("%s %q", server.URL, serverPool.Servers)
	}
	serverPool.DoHealthCheck()
}

func loadBalancing(w http.ResponseWriter, r *http.Request) {
	local := time.Now().Format("15:04:05.000")
	logString := fmt.Sprintf("[%s] %s - %s", local, r.RemoteAddr, r.RequestURI)
	log.Println(logString)

	connection := serverPool.GetNextServer()
	if connection != nil {
		connection.ReverseProxy.ServeHTTP(w, r)
		return
	}

	errorMessage := "Connection Refused. Service unavailable. 502"
	log.Println(errorMessage)
	resp := make(map[string]string)
	resp["message"] = errorMessage
	w.WriteHeader(http.StatusBadGateway)
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	_, writeErr := w.Write(jsonResp)
	if writeErr != nil {
		return
	}
}

func main() {
	cfgData, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = json.Unmarshal(cfgData, &cfg)
	if err != nil {
		return
	}

	createServerPool()

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ProxyPort),
		Handler: http.HandlerFunc(loadBalancing),
	}

	go doHealthCheck()

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
