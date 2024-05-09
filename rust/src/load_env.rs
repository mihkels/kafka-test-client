use std::env;
use dotenv::dotenv;

pub struct EnvVars {
    pub brokers: String,
    pub topic: String,
    pub consumer_group: String,
    pub motivation_filename: String,
    pub lines: usize,
    pub patch_size: usize,
    pub sleep_interval: u64,
    pub run_mode: String
}

pub fn load_env_vars() -> EnvVars {
    dotenv().ok();

    let brokers = env::var("KAFKA_BOOTSTRAP_SERVERS").unwrap_or_else(|_| String::from("localhost:9092"));
    let topic = env::var("TOPIC").unwrap_or_else(|_| String::from("no_motivation"));
    let consumer_group = env::var("CONSUMER_GROUP").unwrap_or_else(|_| String::from("rust-motivation-cg"));
    let motivation_filename = env::var("MOTIVATION_FILE").unwrap_or_else(|_| String::from("motivation.csv"));
    let lines = env::var("NUM_SAMPLES").unwrap_or_else(|_| String::from("10")).parse().unwrap();
    let patch_size = env::var("PATCH_SIZE").unwrap_or_else(|_| String::from("10")).parse().unwrap();
    let sleep_interval = env::var("SLEEP_TIME").unwrap_or_else(|_| String::from("1000")).parse().unwrap();
    let run_mode = env::var("APPLICATION_MODE").unwrap_or_else(|_| String::from("producer")).parse().unwrap();

    EnvVars {
        brokers,
        topic,
        consumer_group,
        motivation_filename,
        lines,
        patch_size,
        sleep_interval,
        run_mode
    }
}