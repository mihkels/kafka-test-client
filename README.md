<div align="center" style="margin: 15px">
    <img with="350" style="padding: 25ps" height="350" src="docs/assets/KafkaKubeClients-logo.svg" />
</div>
<hr />

# KafkaKubeClients - Kafka Client Testing, Tailored for Kubernetes

> We're here to help you check that your Kafka setup works perfectly in Kubernetes, regardless of the programming language you use.

The motivation behind writing the Kafka test client tool was to provide a simple way to test new Kubernetes Kafka cluster deployments and verify that everything was configured correctly. The secondary goal is to provide development teams with configured working samples of producers and consumers in different programming languages.

## Features

- Provide Docker images for all the supported programming languages.
- Supports multiple programming languages: Java, Golang, Rust, Python
- Has a Helm chart that can easily be deployed into the Kubernetes cluster
- The speed of producing data can be configured.
- Has separate image for Kafka producer and consumer
- Provides built-in data set to generate data for producer

## Getting started

The simplest way to begin is by cloning the git repository and deploying the chart from the helm chart directory.

```bash
git clone https://github.com/mihkels/kafka-test-client.git
cd kafka-test-client
# Create namespace to deploy kafka test clients
kubectl create ns kafka-test
# Deploy KafkaKubeClients
helm upgrade kafka-test-client helm-charts/kafka-test-client -n kafka-test -f overrides.yaml --install
```

> 💡 NOTE: You should have a running Kubernetes test cluster with a Kafka cluster installed in your target repository.
> 

Sample `overrides.yaml` file:

```yaml
bootstrapServers: "strimzi-cluster-kafka-bootstrap.kafka.svc:9092"
topic: "motivation"

replicas:
  producer: 1
  consumer: 1

producerConfigurations:
  sampleIntervalInSeconds: 5 # How long to sleep between sending messages
  batchSize: 10 # Number of messages to send in each loop

# The language in which to create the producer and consumer
programmingLanguage:
  - golang
  - python
  - java
  - rust
```

The Kafka test clients have been tested with the following Kafka and Kafka-compliant operators:

- Strimzi Kafka Operator
- Red Panda Operator

Additionally, the following managed Kubernetes offerings have been tested:

- AWS EKS
- Azure AKS
- DigitalOcean Kubernetes cluster
