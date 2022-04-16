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
		serverPool.DoHealthCheck()
	}
}

func loadBalancing(w http.ResponseWriter, r *http.Request) {
	local := time.Now().Format("15:04:05.000")
	_, err := fmt.Fprintf(w, "I'm a load balancer %s", local)
	logString := fmt.Sprintf("[%s] %s - %s", local, r.RemoteAddr, r.RequestURI)
	log.Println(logString)
	if err != nil {
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

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ProxyPort),
		Handler: http.HandlerFunc(loadBalancing),
	}

	go doHealthCheck()

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
