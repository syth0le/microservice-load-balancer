package main

import (
	"encoding/json"
	"fmt"
	"github.com/syth0le/microservice-load-balancer/sources/structures"
	"github.com/syth0le/microservice-load-balancer/sources/utils"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	cfgData, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = json.Unmarshal(cfgData, &structures.Cfg)
	if err != nil {
		return
	}

	utils.CreateServerPool()

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", structures.Cfg.ProxyPort),
		Handler: http.HandlerFunc(structures.LoadBalancing),
	}

	go utils.DoHealthCheck()

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
