package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Statistics struct {
	consumers     map[string]int64
	producers     map[string]int64
	totalSent     int64
	totalReceived int64
}

type WorkerCounter struct {
	Consumers int `json:"consumers"`
	Producers int `json:"producers"`
}

type CollectionRequest struct {
	WorkerType string `json:"workerType"`
	WorkerName string `json:"workerName"`
	Count      int64  `json:"count"`
}

type StatisticsResponse struct {
	TotalConsumers int   `json:"totalConsumers"`
	TotalProducers int   `json:"totalProducers"`
	TotalSent      int64 `json:"totalSent"`
	TotalReceived  int64 `json:"totalReceived"`
}

type CountResponse struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
}

var once sync.Once
var statistics *Statistics

var counterOnce sync.Once
var counterInfo *WorkerCounter

func main() {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/statistics", showStatistics)
	http.HandleFunc("/collect", collectStatistics)
	http.HandleFunc("/worker-count", workerCounter)
	http.Handle("/", jsonMiddleware(http.DefaultServeMux))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting the server: ", err)
		return
	}

	fmt.Println("Server started on port ", port)
}

func workerCounter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	counterOnce.Do(func() {
		counterInfo = &WorkerCounter{
			Consumers: 0,
			Producers: 0,
		}
	})

}

func collectStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	initStatistics()

	var req CollectionRequest
	reqData, _ := getBodyAsString(r)
	fmt.Println(reqData)
	err := json.NewDecoder(strings.NewReader(reqData)).Decode(&req)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	manageStatistics(req)

	data := map[string]string{
		"status": "OK",
	}

	w.WriteHeader(http.StatusOK)
	if convertToJson(w, data) {
		return
	}
}

func initStatistics() {
	once.Do(func() {
		statistics = &Statistics{
			consumers:     make(map[string]int64),
			producers:     make(map[string]int64),
			totalReceived: 0,
			totalSent:     0,
		}
	})
}

func showStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	initStatistics()
	data := StatisticsResponse{
		TotalConsumers: len(statistics.consumers),
		TotalProducers: len(statistics.producers),
		TotalSent:      statistics.totalSent,
		TotalReceived:  statistics.totalReceived,
	}

	if convertToJson(w, data) {
		return
	}

	w.WriteHeader(http.StatusOK)
}

func manageStatistics(req CollectionRequest) {
	if req.WorkerType == "consumer" {
		statistics.totalReceived += req.Count
		statistics.consumers[req.WorkerName] += req.Count
	} else {
		statistics.totalSent += req.Count
		statistics.producers[req.WorkerName] += req.Count
	}
}

func convertToJson(w http.ResponseWriter, d interface{}) bool {
	jsonData, err := json.Marshal(d)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return true
	}

	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return true
	}

	return false
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func getBodyAsString(r *http.Request) (string, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}
