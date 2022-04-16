package main

import (
	"encoding/json"
	"fmt"
	"github.com/syth0le/microservice-load-balancer/config"
	"io/ioutil"
	"log"
)

var cfg config.Config

func doHealthCheck() {
	//for _, server := range config.Servers {
	//if a := server.GetAliveStatus() {
	//	a = 1
	//}
	//pingedURL, err := url.Parse(server.URL)
	//if err != nil {
	//	log.Fatal(err.Error())
	//} else {
	//	server.SetAliveStatus(true)
	//}
	//}
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
	fmt.Println(cfg)
}
