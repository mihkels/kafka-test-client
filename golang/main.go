package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	log.Printf("Starting application\n")
	var appMode = ConfigInstance.ApplicationMode
	log.Printf("Application mode: %s\n", appMode)
	log.Printf("Brokers: %s\n", ConfigInstance.Brokers)
	log.Printf("Topic: %s\n", ConfigInstance.Topic)

	if ConfigInstance.EnableStatistics {
		log.Printf("Statistics collector URL: %s\n", ConfigInstance.StatisticsCollectorURL)
		data, err := json.Marshal(map[string]string{"workerType": appMode})
		if err != nil {
			log.Fatalf("json.Marshal: %s", err)
		}

		client := &http.Client{}
		reader := bytes.NewReader(data)
		resp, err := client.Post(
			ConfigInstance.StatisticsCollectorURL+"/worker-count",
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
		ConfigInstance.WorkerName = result["name"].(string)
	}

	if appMode == "producer" {
		runProducer()
	} else {
		runConsumer()
	}
}
