package main

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"log"
	"math/rand"
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

	pool := compactedTopicUuidPool()
	for {
		for i := 0; i < len(motivationData); i += batchSize {
			end := i + batchSize
			if end > len(motivationData) {
				end = len(motivationData)
			}

			batch := motivationData[i:end]
			ids := []uuid.UUID{}

			for _, row := range batch {
				jsonRow, err := json.Marshal(row)
				if err != nil {
					log.Fatalf("json.Marshal: %s", err)
				}

				var randomUUID uuid.UUID
				if !ConfigInstance.EnableCompactedTopic {
					randomUUID = uuid.New()
				} else {
					randomUUID = pool[rand.Intn(len(pool))]
				}

				ids = append(ids, randomUUID)
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

			if ConfigInstance.EnableStatistics {
				SendStatistics(ConfigInstance.ApplicationMode, ConfigInstance.WorkerName, int64(len(batch)), ids)
			}

			time.Sleep(sleepDuration)
		}
	}
}

func compactedTopicUuidPool() []uuid.UUID {
	var pool = []uuid.UUID{}
	if ConfigInstance.EnableCompactedTopic {
		log.Printf("Compacted topic is enabled with key pool size: %d\n", ConfigInstance.CompactedKeyPoolSize)
		for i := 0; i < ConfigInstance.CompactedKeyPoolSize; i++ {
			poolId := uuid.New()
			pool = append(pool, poolId)
		}
	}
	return pool
}
