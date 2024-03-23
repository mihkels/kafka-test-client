import os

from confluent_kafka import Consumer, KafkaError
from dotenv import load_dotenv

# Load the .env file
load_dotenv()

# Now you can access the variables
kafka_bootstrap_servers = os.getenv('KAFKA_BOOTSTRAP_SERVERS')
motivation_topic = os.getenv('MOTIVATION_TOPIC')
consumer_group_name = os.getenv('MOTIVATION_CONSUMER_GROUP_NAME')

enable_ssl = str(os.getenv('ENABLE_SSL')).lower() == 'true'

# Common Config for Kafka
config = {
    'bootstrap.servers': kafka_bootstrap_servers,
    'group.id': consumer_group_name,
    'auto.offset.reset': 'earliest'
}

# SSL Config for Kafka
if enable_ssl:
    ssl_ca_location = os.getenv('SSL_CA_LOCATION')
    ssl_config = {
        'security.protocol': 'SSL',
        'ssl.ca.location': ssl_ca_location,
    }
    config.update(ssl_config)

# Create a Kafka consumer
consumer = Consumer(config)

# Subscribe to the topic
consumer.subscribe([motivation_topic])

while True:
    msg = consumer.poll(1.0)

    if msg is None:
        continue
    if msg.error():
        if msg.error().code() == KafkaError.get_partition_eof():
            continue
        else:
            print(msg.error())
            break

    print('Received message: {}'.format(msg.value().decode('utf-8')))

consumer.close()