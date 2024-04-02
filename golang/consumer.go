package main

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type ConsumerGroupHandler struct{}

func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("Message topic:%q, body:%q, partition:%d offset:%d\n", msg.Topic, string(msg.Value), msg.Partition, msg.Offset)
		sess.MarkMessage(msg, "")
	}
	return nil
}

func runConsumer() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	group, err := sarama.NewConsumerGroup([]string{ConfigInstance.Brokers}, ConfigInstance.Group, config)
	if err != nil {
		log.Fatalln("Error creating consumer group client: ", err)
	}

	ctx := context.Background()
	handler := ConsumerGroupHandler{}

	go func() {
		for {
			if err := group.Consume(ctx, []string{ConfigInstance.Topic}, handler); err != nil {
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