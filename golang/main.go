package main

import (
	"log"
)

func main() {
	log.Printf("Starting application\n")
	var appMode = ConfigInstance.ApplicationMode
	log.Printf("Application mode: %s\n", appMode)
	log.Printf("Brokers: %s\n", ConfigInstance.Brokers)
	log.Printf("Topic: %s\n", ConfigInstance.Topic)
	if appMode == "producer" {
		runProducer()
	} else {
		runConsumer()
	}
}
