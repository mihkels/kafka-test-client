use futures::stream::StreamExt;
use rdkafka::config::ClientConfig;
use rdkafka::consumer::{Consumer, StreamConsumer};
use rdkafka::Message;
use crate::load_env::load_env_vars;

pub async fn run_consumer() {
    let env_vars = load_env_vars();

    let consumer: StreamConsumer = ClientConfig::new()
        .set("bootstrap.servers", &env_vars.brokers)
        .set("group.id", &env_vars.consumer_group)
        .create()
        .expect("Consumer creation failed");

    consumer
        .subscribe(&[&env_vars.topic])
        .expect("Can't subscribe to specified topic");

    let mut message_stream = consumer.stream();

    while let Some(message) = message_stream.next().await {
        match message {
            Ok(msg) => {
                let payload = msg.payload_view::<str>().unwrap().unwrap();
                println!("Received Message: {}", payload);
            },
            Err(err) => println!("Error: {:?}", err),
        }
    }
}