package sources

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	Attempts int = iota
	Retry
)

func LoadBalancing(w http.ResponseWriter, r *http.Request) {
	attempts := GetAttemptsFromContext(r)
	if attempts > 3 {
		log.Printf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}

	local := time.Now().Format("15:04:05.000")
	logString := fmt.Sprintf("[%s] %s - %s", local, r.RemoteAddr, r.RequestURI)
	log.Println(logString)

	connection := ServPool.GetNextServer()
	if connection != nil {
		connection.ReverseProxy.ServeHTTP(w, r)
		return
	}

	errorMessage := "Connection Refused. Service unavailable. 502"
	log.Println(errorMessage)
	resp := make(map[string]string)
	resp["message"] = errorMessage
	w.WriteHeader(http.StatusServiceUnavailable)
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	_, writeErr := w.Write(jsonResp)
	if writeErr != nil {
		return
	}
}

func GetAttemptsFromContext(r *http.Request) int {
	if attempts, ok := r.Context().Value(Attempts).(int); ok {
		return attempts
	}
	return 1
}

func GetRetryFromContext(r *http.Request) int {
	if retry, ok := r.Context().Value(Retry).(int); ok {
		return retry
	}
	return 0
}
