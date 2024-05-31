# kafka-test-client

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

> 💡 NOTE: You should have a running Kubernetes test cluster with a Kafka cluster installed in your target repository.
> 

The Kafka test clients have been tested with the following Kafka and Kafka-compliant operators:

- Strimzi Kafka Operator
- Red Panda Operator

Additionally, the following managed Kubernetes offerings have been tested:

- AWS EKS
- Azure AKS
- DigitalOcean Kubernetes cluster
