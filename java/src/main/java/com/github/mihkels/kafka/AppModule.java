package com.github.mihkels.kafka;

import com.fasterxml.jackson.databind.ObjectMapper;
import dagger.Module;
import dagger.Provides;
import io.github.mihkels.dotenv.Dotenv;
import org.apache.kafka.clients.producer.KafkaProducer;
import org.apache.kafka.clients.producer.Producer;

import java.util.Objects;
import java.util.Properties;

@Module
public class AppModule {
    @Provides
    public ObjectMapper objectMapper() {
        return new ObjectMapper();
    }

    @Provides
    public CSVReader provideCSVReader() {
        return new CSVReader();
    }

    @Provides
    public AppConfig provideAppConfig() {
        var dotenv = Dotenv.configure().ignoreIfMissing().load();

        // Kafka producer configuration
        var keySerializer = "org.apache.kafka.common.serialization.StringSerializer";
        var valueSerializer = "org.apache.kafka.common.serialization.StringSerializer";
        var bootstrapServers = dotenv.get("KAFKA_BOOTSTRAP_SERVERS");
        var topicName = dotenv.get("TOPIC");
        var csvFilePath = dotenv.get("MOTIVATION_FILE");
        var numberOfLines = Integer.parseInt(Objects.requireNonNull(dotenv.get("NUM_SAMPLES", "0")));
        var patchSize = Integer.parseInt(Objects.requireNonNull(dotenv.get("BATCH_SIZE", "0")));
        var sleepTime = Integer.parseInt(Objects.requireNonNull(dotenv.get("SLEEP_TIME", "0")));
        var consumerGroupId = dotenv.get("MOTIVATION_CONSUMER_GROUP_NAME");
        var applicationType = dotenv.get("APPLICATION_MODE");

        return new AppConfig(
                bootstrapServers,
                topicName,
                keySerializer,
                valueSerializer,
                csvFilePath,
                consumerGroupId,
                applicationType,
                numberOfLines,
                patchSize,
                sleepTime,
                true
        );
    }

    @Provides
    public DataLoaderConfig provideDataLoaderConfig(AppConfig appConfig) {
        return new DataLoaderConfig(
                appConfig.csvFilePath(),
                appConfig.numberOfLines(),
                appConfig.patchSize(),
                appConfig.sleepTime()
        );
    }

    @Provides
    public Producer<String, String> provideKafkaProducer(AppConfig config) {
        Properties props = new Properties();
        props.put("bootstrap.servers", config.bootstrapServers());
        props.put("key.serializer", config.keySerializer());
        props.put("value.serializer", config.valueSerializer());

        return new KafkaProducer<>(props);
    }
}
