package main

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
	"time"
)

func runProducer() {
	motivationData := readMotivationFile()
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{ConfigInstance.Brokers}, config)
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	batchSize := ConfigInstance.BatchSize
	sleepDuration := time.Second * ConfigInstance.SleepInterval

	for {
		for i := 0; i < len(motivationData); i += batchSize {
			end := i + batchSize
			if end > len(motivationData) {
				end = len(motivationData)
			}

			batch := motivationData[i:end]

			for _, row := range batch {
				jsonRow, err := json.Marshal(row)
				if err != nil {
					log.Fatalf("json.Marshal: %s", err)
				}

				msg := &sarama.ProducerMessage{
					Topic: ConfigInstance.Topic,
					Value: sarama.StringEncoder(string(jsonRow)),
				}
				partition, offset, err := producer.SendMessage(msg)
				if err != nil {
					log.Fatalln(err)
				}
				log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", ConfigInstance.Topic, partition, offset)
			}

			time.Sleep(sleepDuration)
		}
	}
}
