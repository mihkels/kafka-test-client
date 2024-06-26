package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Statistics struct {
	consumers         map[string]int64
	producers         map[string]int64
	totalSent         int64
	totalReceived     int64
	matchedIds        int64
	missingMatchedIds int64
}

type WorkerCounter struct {
	Consumers           int      `json:"consumers"`
	Producers           int      `json:"producers"`
	ConsumerWorkerNames []string `json:"consumerWorkerNames"`
	ProducerWorkerNames []string `json:"producerWorkerNames"`
}

type CollectionRequest struct {
	WorkerType string      `json:"workerType"`
	WorkerName string      `json:"workerName"`
	Count      int64       `json:"count"`
	Ids        []uuid.UUID `json:"ids"`
}

type StatisticsResponse struct {
	TotalConsumers int   `json:"totalConsumers"`
	TotalProducers int   `json:"totalProducers"`
	TotalSent      int64 `json:"totalSent"`
	TotalReceived  int64 `json:"totalReceived"`
	MissingIds     int64 `json:"missingIds"`
	MatchedIds     int64 `json:"matchedIds"`
}

type CountResponse struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
}

var counterMutex = &sync.Mutex{}
var statisticsMutex = &sync.Mutex{}

var statistics *Statistics

var counterInfo *WorkerCounter
var idMatcher map[uuid.UUID]bool

func main() {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	initWorkerCounter()
	initStatistics()

	http.HandleFunc("/statistics", jsonMiddleware(showStatistics))
	http.HandleFunc("/statistics/reset", jsonMiddleware(resetStatistics))
	http.HandleFunc("/collect", jsonMiddleware(collectStatistics))
	http.HandleFunc("/worker-count", jsonMiddleware(workerCounter))
	http.HandleFunc("/worker-count/info", jsonMiddleware(workerCounterInfo))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting the server: ", err)
		return
	}

	fmt.Println("Server started on port ", port)
}

func workerCounterInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	if convertToJson(w, counterInfo) {
		return
	}
}

func initWorkerCounter() {
	counterInfo = &WorkerCounter{
		Consumers:           0,
		Producers:           0,
		ConsumerWorkerNames: make([]string, 0),
		ProducerWorkerNames: make([]string, 0),
	}
}

func workerCounter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var reqMap map[string]interface{}
	err = json.Unmarshal(bodyBytes, &reqMap)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	workerType, ok := reqMap["workerType"].(string)
	if !ok {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var resp CountResponse
	counterMutex.Lock()
	if workerType == "consumer" {
		consumerName := "consumer-" + fmt.Sprint(counterInfo.Consumers)
		counterInfo.ConsumerWorkerNames = append(counterInfo.ConsumerWorkerNames, consumerName)
		counterInfo.Consumers += 1
		resp = CountResponse{
			Count: counterInfo.Consumers,
			Name:  consumerName,
		}
	} else {
		producerName := "producer-" + fmt.Sprint(counterInfo.Producers)
		counterInfo.Producers += 1
		counterInfo.ProducerWorkerNames = append(counterInfo.ProducerWorkerNames, producerName)

		resp = CountResponse{
			Count: counterInfo.Producers,
			Name:  producerName,
		}
	}
	counterMutex.Unlock()

	w.WriteHeader(http.StatusOK)
	if convertToJson(w, resp) {
		return
	}
}

func collectStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CollectionRequest
	reqData, _ := getBodyAsString(r)
	fmt.Println(reqData)
	err := json.NewDecoder(strings.NewReader(reqData)).Decode(&req)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	manageStatistics(req)
	handleMatchedIds(req)

	data := map[string]string{
		"status": "OK",
	}

	w.WriteHeader(http.StatusOK)
	if convertToJson(w, data) {
		return
	}
}

func initStatistics() {
	statistics = &Statistics{
		consumers:         make(map[string]int64),
		producers:         make(map[string]int64),
		totalReceived:     0,
		totalSent:         0,
		missingMatchedIds: 0,
		matchedIds:        0,
	}
}

func resetStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	initStatistics()
	initWorkerCounter()

	w.WriteHeader(http.StatusOK)
	if convertToJson(w, map[string]string{"status": "OK"}) {
		return
	}
}

func showStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := StatisticsResponse{
		TotalConsumers: len(statistics.consumers),
		TotalProducers: len(statistics.producers),
		TotalSent:      statistics.totalSent,
		TotalReceived:  statistics.totalReceived,
		MissingIds:     statistics.missingMatchedIds,
		MatchedIds:     statistics.matchedIds,
	}

	w.WriteHeader(http.StatusOK)
	if convertToJson(w, data) {
		return
	}
}

func manageStatistics(req CollectionRequest) {
	statisticsMutex.Lock()
	defer statisticsMutex.Unlock()
	if req.WorkerType == "consumer" {
		statistics.totalReceived += req.Count
		statistics.consumers[req.WorkerName] += req.Count
	}

	if req.WorkerType == "producer" {
		statistics.totalSent += req.Count
		statistics.producers[req.WorkerName] += req.Count
	}
}

func handleMatchedIds(req CollectionRequest) {
	statisticsMutex.Lock()
	defer statisticsMutex.Unlock()

	if idMatcher == nil {
		idMatcher = make(map[uuid.UUID]bool)
	}

	if req.WorkerType == "consumer" {
		for _, id := range req.Ids {
			if _, ok := idMatcher[id]; ok {
				statistics.matchedIds++
				delete(idMatcher, id)
			}
		}
	}

	if req.WorkerType == "producer" {
		for _, id := range req.Ids {
			idMatcher[id] = true
		}
	}

	statistics.missingMatchedIds = int64(len(idMatcher))
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

func jsonMiddleware(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func getBodyAsString(r *http.Request) (string, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}
