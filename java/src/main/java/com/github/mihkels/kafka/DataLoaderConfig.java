package com.github.mihkels.kafka;

public record DataLoaderConfig(
        String csvFilePath,
        int numberOfLines,
        int patchSize,
        int sleepTime
) {}
