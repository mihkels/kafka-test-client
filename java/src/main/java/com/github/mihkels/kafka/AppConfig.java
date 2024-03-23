package com.github.mihkels.kafka;

public record AppConfig(
        String bootstrapServers,
        String topicName,
        String keySerializer,
        String valueSerializer,
        String csvFilePath,
        String groupId,
        String applicationType,
        int numberOfLines,
        int patchSize,
        int sleepTime,
        boolean stopApplication
) {}
