package com.github.mihkels.kafka;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.inject.Inject;
import java.util.Collections;
import java.util.List;
import java.util.Map;

public class MotivationDataLoader {
    public static final Logger logger = LoggerFactory.getLogger(MotivationDataLoader.class);
    private final DataLoaderConfig config;
    private final CSVReader csvReader;

    @Inject
    public MotivationDataLoader(DataLoaderConfig config, CSVReader csvReader) {
        this.config = config;
        this.csvReader = csvReader;
    }

    public List<Map<String, String>> loadMotivationData() {
        logger.info("Loading motivation data");
        var data = csvReader.readCSV(config.csvFilePath());
        logger.info("Loaded {} rows", data.size());
        Collections.shuffle(data);
        return data.subList(0, Math.min(config.numberOfLines(), data.size()));
    }
}

