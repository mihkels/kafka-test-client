package com.github.mihkels.kafka;

import org.apache.kafka.clients.consumer.Consumer;
import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.clients.consumer.ConsumerRecords;
import org.apache.kafka.clients.consumer.KafkaConsumer;
import org.apache.kafka.common.serialization.StringDeserializer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.inject.Inject;
import java.time.Duration;
import java.util.Collections;
import java.util.Properties;

public class MotivationKafkaConsumer {
    private static final Logger logger = LoggerFactory.getLogger(MotivationKafkaConsumer.class);

    private final AppConfig config;
    private final Consumer<String, String> consumer;

    @Inject
    public MotivationKafkaConsumer(AppConfig config) {
        this.config = config;
        Properties props = new Properties();
        props.put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG, config.bootstrapServers());
        props.put(ConsumerConfig.GROUP_ID_CONFIG, config.groupId());
        props.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class.getName());
        props.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class.getName());
        this.consumer = new KafkaConsumer<>(props);
    }

    public void run() {
        consumer.subscribe(Collections.singletonList(config.topicName()));

        while (config.stopApplication()) {
            ConsumerRecords<String, String> records = consumer.poll(Duration.ofMillis(1000));
            records.forEach(message -> logger.info("Received message: topic = {}, partition = {}, offset = {}, key = {}, value = {}",
                    message.topic(), message.partition(), message.offset(), message.key(), message.value()));
        }
    }
}
