package com.github.mihkels.kafka;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.kafka.clients.producer.Producer;
import org.apache.kafka.clients.producer.ProducerRecord;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.inject.Inject;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

public class MotivationKafkaProducer {
    public static final Logger logger = LoggerFactory.getLogger(MotivationKafkaProducer.class);

    private final AppConfig config;
    private final MotivationDataLoader motivationDataLoader;
    private final ObjectMapper objectMapper;
    private final Producer<String, String> producer;

    @Inject
    public MotivationKafkaProducer(
            Producer<String, String > producer,
            AppConfig config,
            MotivationDataLoader dataLoader,
            ObjectMapper objectMapper
    ) {
        this.producer = producer;
        this.config = config;
        this.motivationDataLoader = dataLoader;
        this.objectMapper = objectMapper;
    }

    public void run() {
        logger.info("Starting the producer");
        var sendData = loadMotivationData();
        while (config.stopApplication()) {
            Collections.shuffle(sendData);
            List<String> randomLines = sendData.subList(0, Math.min(config.patchSize(), sendData.size()));

            randomLines.forEach(data -> {
                var message = new ProducerRecord<String, String>(config.topicName(), data);
                producer.send(message, (metadata, exception) -> {
                    if (exception == null) {
                        logger.info("Message sent to topic {}, partition {}, offset {}", metadata.topic(), metadata.partition(), metadata.offset());
                    } else {
                        logger.warn("Failed to send message", exception);
                    }
                });
            });

            try {
                // Sleep for a given number of seconds
                Thread.sleep(config.sleepTime() * 1000L);
            } catch (InterruptedException e) {
                logger.error("Producer interrupted", e);
                // Restore interrupted state...
                Thread.currentThread().interrupt();
            }
        }

        // Close the producer
        producer.close();
        logger.info("Producer finished");
    }

    private List<String> loadMotivationData() {
        logger.info("Loading motivation data");
        var data = motivationDataLoader.loadMotivationData()
                .stream()
                .map(this::convertMapToJson)
                .collect(Collectors.toCollection(ArrayList::new));

        logger.info("Loaded {} rows", data.size());
        return data;
    }

    private String convertMapToJson(Map<String, String> map) {
        try {
            return objectMapper.writeValueAsString(map);
        } catch (JsonProcessingException e) {
            logger.error("Error converting map to JSON", e);
            return null;
        }
    }
}
