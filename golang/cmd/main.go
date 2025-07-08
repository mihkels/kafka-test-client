package main

import (
	"bytes"
	"encoding/json"
	"kafka_demo/internal/client"
	"kafka_demo/internal/config"
	"log"
	"net/http"
)

func main() {
	log.Printf("Starting application\n")
	var appMode = config.Instance.ApplicationMode
	log.Printf("Application mode: %s\n", appMode)
	log.Printf("Brokers: %s\n", config.Instance.Brokers)
	log.Printf("Topic: %s\n", config.Instance.Topic)

	if config.Instance.EnableStatistics {
		log.Printf("Statistics collector URL: %s\n", config.Instance.StatisticsCollectorURL)
		data, err := json.Marshal(map[string]string{"workerType": appMode})
		if err != nil {
			log.Fatalf("json.Marshal: %s", err)
		}

		statClient := &http.Client{}
		reader := bytes.NewReader(data)
		resp, err := statClient.Post(
			config.Instance.StatisticsCollectorURL+"/worker-count",
			"application/json",
			reader)
		if err != nil {
			log.Fatalf("client.Post: %s", err)
		}

		defer resp.Body.Close()
		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			log.Fatalf("json.NewDecoder: %s", err)
		}

		log.Printf("Worker count: %v\n", result)
		config.Instance.WorkerName = result["name"].(string)
	}

	if appMode == "producer" {
		client.RunProducer()
	} else {
		client.RunConsumer()
	}
}
