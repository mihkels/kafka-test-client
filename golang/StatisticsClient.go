package main

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
)

type CollectionRequest struct {
	WorkerType string      `json:"workerType"`
	WorkerName string      `json:"workerName"`
	Count      int64       `json:"count"`
	Ids        []uuid.UUID `json:"ids"`
}

func SendStatistics(workerType string, workerName string, count int64, ids []uuid.UUID) {
	url := ConfigInstance.StatisticsCollectorURL
	req := CollectionRequest{
		WorkerType: workerType,
		WorkerName: workerName,
		Count:      count,
		Ids:        ids,
	}

	data, err := json.Marshal(req)
	if err != nil {
		log.Printf("json.Marshal failed: %s", err)
		return
	}

	resp, err := http.Post(url+"/collect", "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("http.Post failed: %s", err)
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	log.Println("Statistics sent to the statistics collector", resp.Status)
}
