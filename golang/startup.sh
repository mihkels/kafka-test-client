#!/usr/bin/env sh

echo "Starting Kafka test client"
echo "Brokers: ${KAFKA_BOOTSTRAP_SERVERS}"
echo "Topic: ${TOPIC}"

/app/kafka_app