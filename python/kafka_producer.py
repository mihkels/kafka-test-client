import json
import os
import time

import pandas as pd
from confluent_kafka import Producer
from dotenv import load_dotenv

load_dotenv()

num_samples = int(os.getenv('NUM_SAMPLES', 100))
batch_size = int(os.getenv('BATCH_SIZE', 10))
sleep_time = int(os.getenv('SLEEP_TIME', 5))
motivation_file = os.getenv('MOTIVATION_FILE')

kafka_bootstrap_servers = os.getenv('KAFKA_BOOTSTRAP_SERVERS')
motivation_topic = os.getenv('MOTIVATION_TOPIC')

enable_ssl = str(os.getenv('ENABLE_SSL')).lower() == 'true'


# Common Config for Kafka
config = {
    'bootstrap.servers': kafka_bootstrap_servers,
}

# SSL Config for Kafka
if enable_ssl:
    print("Enabling SSL")
    ssl_ca_location = os.getenv('SSL_CA_LOCATION')
    ssl_config = {
        'security.protocol': 'SSL',
        'ssl.ca.location': ssl_ca_location,
    }
    config |= ssl_config

producer = Producer(config)

df = pd.read_csv(motivation_file)
random_rows = df.sample(n=num_samples)


def delivery_report(err, msg):
    """ Called once for each message produced to indicate delivery result.
        Triggered by poll() or flush(). """
    if err is not None:
        print(f'Message delivery failed: {err}')
    else:
        print(f'Message delivered to {msg.topic()} [{msg.partition()}]')


while True:
    try:
        for i in range(0, len(random_rows), batch_size):
            batch = random_rows[i:i + batch_size]
            for index, row in batch.iterrows():
                # Convert row to dictionary

                row_dict = {k.lower(): v for k, v in row.to_dict().items()}
                print(row_dict)
                # Send the row to Kafka topic
                producer.produce(motivation_topic, json.dumps(row_dict).encode('utf-8'), callback=delivery_report)

            # Trigger any available delivery report callbacks from previous produce() calls
            producer.poll(0)

            # Sleep for 5 seconds
            time.sleep(sleep_time)
    except KeyboardInterrupt:
        print("\nInterrupted. Flushing pending messages...")
    finally:
        # Wait for any outstanding messages to be delivered and delivery reports to be received.
        producer.flush()
