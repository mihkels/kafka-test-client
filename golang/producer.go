package main

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
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

	log.Printf("Start sending data into topic: %s\n", ConfigInstance.Topic)
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

				randomUUID := uuid.New()
				msg := &sarama.ProducerMessage{
					Topic: ConfigInstance.Topic,
					Value: sarama.StringEncoder(jsonRow),
					Key:   sarama.StringEncoder(randomUUID.String()),
				}
				partition, offset, err := producer.SendMessage(msg)
				if err != nil {
					log.Fatalln(err)
				}
				log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d) with key: %s\n", ConfigInstance.Topic, partition, offset, msg.Key)
			}

			SendStatistics(ConfigInstance.ApplicationMode, ConfigInstance.WorkerName, int64(len(batch)))
			time.Sleep(sleepDuration)
		}
	}
}
