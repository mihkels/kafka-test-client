use std::fs::File;
use std::time::Duration;
use csv::{ReaderBuilder};
use rand::seq::IteratorRandom;
use serde_json::{json, Value};
use crate::load_env::load_env_vars;
use rdkafka::config::ClientConfig;
use rdkafka::producer::{FutureProducer, FutureRecord};
use rdkafka::util::Timeout;

pub fn load_csv_and_select_lines(filename: &str, lines: usize) -> Vec<Value> {
    let file = File::open(filename).unwrap();
    let mut reader = ReaderBuilder::new().has_headers(true).from_reader(file);
    let headers = reader.headers().unwrap().clone();
    let records: Vec<_> = reader.records().collect::<Result<_, _>>().unwrap();
    let mut rng = rand::thread_rng();
    let random_records: Vec<_> = records.iter().choose_multiple(&mut rng, lines);

    let mut json_records= Vec::new();
    for record in random_records {
        let mut json_record = json!({});
        for (header, field) in headers.iter().zip(record.iter()) {
            if header == "Tags" {
                let tags_list: Vec<&str> = field.split(',').collect();
                json_record[header.to_lowercase()] = json!(tags_list);
            } else {
                json_record[header.to_lowercase()] = json!(field);
            }
        }
        json_records.push(json_record);
    }

    json_records
}
pub async fn run_producer() {
    let env_vars = load_env_vars();

    let data = load_csv_and_select_lines(&env_vars.motivation_filename, env_vars.lines);
    let batch_size = env_vars.patch_size;
    let sleep_interval = Duration::from_secs(env_vars.sleep_interval);

    let batches: Vec<_> = data.chunks(batch_size).collect();

    let producer: FutureProducer = ClientConfig::new()
        .set("bootstrap.servers", &env_vars.brokers)
        .set("message.timeout.ms", "5000")
        .create()
        .expect("Producer creation error");

    loop {
        for batch in &batches {
            for record in *batch {
                let message = record.to_string();
                let delivery_future = producer.send(
                    FutureRecord::to(&env_vars.topic)
                        .payload(&message)
                        .key("motivation"),
                    Timeout::Never,
                );

                match delivery_future.await {
                    Ok((partition, offset)) => println!(
                        "Message delivered to partition: {} and offset: {}",
                        partition, offset
                    ),
                    Err((err, _)) => println!(
                        "Failed to deliver message; reason: {}",
                        err
                    ),
                }
            }
            tokio::time::sleep(sleep_interval).await;
        }
    }
}