package client

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"kafka_demo/internal/config"
	"kafka_demo/internal/stats"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type ConsumerGroupHandler struct {
	messageCount int
	ids          []uuid.UUID
	mu           sync.Mutex
}

func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("Message topic:%q, body:%q, partition:%d offset:%d\n", msg.Topic, string(msg.Value), msg.Partition, msg.Offset)
		sess.MarkMessage(msg, "")

		h.mu.Lock()
		h.messageCount++
		id, _ := uuid.Parse(string(msg.Key))
		h.ids = append(h.ids, id)
		if h.messageCount == 10 {
			if config.Instance.EnableStatistics {
				stats.SendStatistics(config.Instance.ApplicationMode, config.Instance.WorkerName, int64(h.messageCount), h.ids)
			}
			h.messageCount = 0
			h.ids = []uuid.UUID{}
		}
		h.mu.Unlock()
	}
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	h.mu.Lock()
	if h.messageCount > 0 {
		if config.Instance.EnableStatistics {
			stats.SendStatistics(config.Instance.ApplicationMode, config.Instance.WorkerName, int64(h.messageCount), h.ids)
		}

		h.messageCount = 0
		h.ids = []uuid.UUID{}
	}
	h.mu.Unlock()
	return nil
}

func RunConsumer() {
	clientConfig := KafkaClient{
		Config: config.Instance,
	}

	kafkaConfig, err := clientConfig.NewKafkaClientConfig()
	kafkaConfig.Consumer.Return.Errors = true
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	group, err := sarama.NewConsumerGroup([]string{config.Instance.Brokers}, config.Instance.Group, kafkaConfig)
	if err != nil {
		log.Fatalln("Error creating consumer group client: ", err)
	}

	ctx := context.Background()
	handler := ConsumerGroupHandler{}

	go func() {
		for {
			if err := group.Consume(ctx, []string{config.Instance.Topic}, &handler); err != nil {
				log.Fatalln("Error from consumer: ", err)
			}
		}
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm

	if err := group.Close(); err != nil {
		log.Fatalln("Error closing client: ", err)
	}
}
