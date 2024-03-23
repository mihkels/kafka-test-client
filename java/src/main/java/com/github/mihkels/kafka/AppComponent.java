package com.github.mihkels.kafka;
import dagger.Component;

@Component(modules = AppModule.class)
public interface AppComponent {
    MotivationKafkaProducer producer();
    MotivationKafkaConsumer consumer();
    MotivationDataLoader dataLoader();
    AppConfig provideAppConfig(); // Add this line
}
