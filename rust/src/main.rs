
mod consumer;
mod producer;
mod load_env;

#[tokio::main]
async fn main() {
    load_env::load_env_vars();
    if load_env::load_env_vars().run_mode == "consumer" {
        consumer::run_consumer().await;
    } else {
        producer::run_producer().await;
    }
}