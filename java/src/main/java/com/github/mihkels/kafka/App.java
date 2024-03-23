package com.github.mihkels.kafka;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class App {
    private static final Logger logger = LoggerFactory.getLogger(App.class);

    public static void main( String[] args) {
        logger.info("Starting the application");

        var appComponent = DaggerAppComponent.create();
        var config = appComponent.provideAppConfig();
        if (config.applicationType().equals("producer")) {
            var producer = appComponent.producer();
            producer.run();
        } else {
            var consumer = appComponent.consumer();
            consumer.run();
        }

        logger.info("Application finished");
    }
}
