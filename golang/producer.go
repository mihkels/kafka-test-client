package main

import (
	"crypto/sha256"
	"encoding/hex"
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
			ids := []uuid.UUID{}

			for i, row := range batch {
				jsonRow, err := json.Marshal(row)
				if err != nil {
					log.Fatalf("json.Marshal: %s", err)
				}

				randomUUID := uuid.New()
				ids = append(ids, randomUUID)
				msg := &sarama.ProducerMessage{
					Topic: ConfigInstance.Topic,
					Value: sarama.StringEncoder(jsonRow),
					Key:   sarama.StringEncoder(randomUUID.String()),
				}

				if ConfigInstance.UseHeaders {
					msg.Headers = addHeaders(int64(i), string(jsonRow))
				}

				partition, offset, err := producer.SendMessage(msg)
				if err != nil {
					log.Fatalln(err)
				}
				log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d) with key: %s\n", ConfigInstance.Topic, partition, offset, msg.Key)
			}

			SendStatistics(ConfigInstance.ApplicationMode, ConfigInstance.WorkerName, int64(len(batch)), ids)
			time.Sleep(sleepDuration)
		}
	}
}

func addHeaders(counter int64, data string) []sarama.RecordHeader {
	return []sarama.RecordHeader{
		{
			Key:   []byte("hash"),
			Value: []byte(getSHA256Hash(string(data))),
		},
		{
			Key:   []byte("time"),
			Value: []byte(time.Now().Format(time.RFC3339)),
		},
		{
			Key:   []byte("position"),
			Value: []byte(string(counter)),
		},
		{
			Key:   []byte("producer"),
			Value: []byte(ConfigInstance.ApplicationMode + "-" + ConfigInstance.WorkerName),
		},
	}
}

func getSHA256Hash(text string) string {
	hash := sha256.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}
