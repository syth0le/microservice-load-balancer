package balancer

import (
	"encoding/json"
	"fmt"
	"github.com/syth0le/microservice-load-balancer/config"
	"log"
	"net/http"
	"time"
)

func LoadBalancing(w http.ResponseWriter, r *http.Request) {
	local := time.Now().Format("15:04:05.000")
	logString := fmt.Sprintf("[%s] %s - %s", local, r.RemoteAddr, r.RequestURI)
	log.Println(logString)

	connection := config.ServerPool.GetNextServer()
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
