package client

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	config2 "kafka_demo/internal/config"
	"kafka_demo/internal/stats"
	"log"
	"math/rand"
	"time"
)

func RunProducer() {
	motivationData := readMotivationFile()
	clientConfig := KafkaClient{
		Config: config2.Instance,
	}

	config, err := clientConfig.NewKafkaClientConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{config2.Instance.Brokers}, config)
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	batchSize := config2.Instance.BatchSize
	sleepDuration := time.Second * config2.Instance.SleepInterval

	log.Printf("Start sending data into topic: %s\n", config2.Instance.Topic)

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
				if !config2.Instance.EnableCompactedTopic {
					randomUUID = uuid.New()
				} else {
					randomUUID = pool[rand.Intn(len(pool))]
				}

				ids = append(ids, randomUUID)
				msg := &sarama.ProducerMessage{
					Topic: config2.Instance.Topic,
					Value: sarama.StringEncoder(jsonRow),
					Key:   sarama.StringEncoder(randomUUID.String()),
				}
				partition, offset, err := producer.SendMessage(msg)
				if err != nil {
					log.Fatalln(err)
				}
				log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d) with key: %s\n", config2.Instance.Topic, partition, offset, msg.Key)
			}

			if config2.Instance.EnableStatistics {
				stats.SendStatistics(config2.Instance.ApplicationMode, config2.Instance.WorkerName, int64(len(batch)), ids)
			}

			time.Sleep(sleepDuration)
		}
	}
}

func compactedTopicUuidPool() []uuid.UUID {
	var pool = []uuid.UUID{}
	if config2.Instance.EnableCompactedTopic {
		log.Printf("Compacted topic is enabled with key pool size: %d\n", config2.Instance.CompactedKeyPoolSize)
		for i := 0; i < config2.Instance.CompactedKeyPoolSize; i++ {
			poolId := uuid.New()
			pool = append(pool, poolId)
		}
	}
	return pool
}
