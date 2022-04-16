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
		countdown := time.NewTicker(1 * time.Minute)
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
		serverPool.AddServer(&server)
	}
}

func loadBalancing(w http.ResponseWriter, r *http.Request) {
	local := time.Now().Format("15:04:05.000")
	//_, err := fmt.Fprintf(w, "I'm a load balancer %s", local)
	logString := fmt.Sprintf("[%s] %s - %s", local, r.RemoteAddr, r.RequestURI)
	log.Println(logString)
	//if err != nil {
	//	return
	//}

	connection := serverPool.GetNextServer()
	for _, server := range serverPool.Servers {
		log.Println(server, serverPool.CurrentServer)
	}
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
	w.Write(jsonResp)

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
