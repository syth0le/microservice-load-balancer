package main

import (
	"encoding/json"
	"fmt"
	"github.com/syth0le/microservice-load-balancer/sources"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	cfgData, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = json.Unmarshal(cfgData, &sources.Cfg)
	if err != nil {
		return
	}

	sources.CreateServerPool()

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", sources.Cfg.ProxyPort),
		Handler: http.HandlerFunc(sources.LoadBalancing),
	}

	go sources.DoHealthCheck()

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
