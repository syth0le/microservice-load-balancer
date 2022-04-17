package main

import (
	"encoding/json"
	"fmt"
	"github.com/syth0le/microservice-load-balancer/config"
	"github.com/syth0le/microservice-load-balancer/sources/balancer"
	"github.com/syth0le/microservice-load-balancer/sources/structures"
	"github.com/syth0le/microservice-load-balancer/sources/utils"
	"io/ioutil"
	"log"
	"net/http"
)

var ServerPool structures.ServerPool

func main() {
	cfgData, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = json.Unmarshal(cfgData, &config.Cfg)
	if err != nil {
		return
	}

	utils.CreateServerPool()

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", config.Cfg.ProxyPort),
		Handler: http.HandlerFunc(balancer.LoadBalancing),
	}

	go utils.DoHealthCheck()

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
