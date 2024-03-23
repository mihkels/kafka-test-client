package com.github.mihkels.kafka;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.io.FileNotFoundException;
import java.util.*;

public class CSVReader {
    private static final Logger logger = LoggerFactory.getLogger(CSVReader.class);

    public List<Map<String, String>> readCSV(String filePath) {
        logger.info("Reading CSVj file: {}", filePath);
        List<Map<String, String>> data = new ArrayList<>();
        try (Scanner scanner = new Scanner(new File(filePath))) {
            if (!scanner.hasNext()) {
                return data;
            }

            String[] headers = readHeaders(scanner);
            readValuesIntoMap(scanner, headers, data);
        } catch (FileNotFoundException e) {
            logger.error("File not found", e);
        }

        return data;
    }

    private static void readValuesIntoMap(Scanner scanner, String[] headers, List<Map<String, String>> data) {
        while (scanner.hasNext()) {
            String[] values = scanner.nextLine().split(",");
            Map<String, String> row = new HashMap<>();
            for (int i = 0; i < headers.length; i++) {
                if (i < values.length) {
                    row.put(headers[i], values[i]);
                } else {
                    row.put(headers[i], ""); // or skip this header
                }
            }
            logger.debug("Read row: {}", row);
            data.add(row);
        }
    }

    private static String [] readHeaders(Scanner scanner) {
        String[] headers = scanner.nextLine().split(",");
        for (int i = 0; i < headers.length; i++) {
            headers[i] = headers[i].trim().toLowerCase();
        }
        return headers;
    }
}
