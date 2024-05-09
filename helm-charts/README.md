```shell
 helm upgrade --install kafka-test-client . -n kafka-worker -f overrides.yaml --dry-run
```

Create topic inside the Kafka cluster:

```shell
bin/kafka-topics.sh --bootstrap-server localhost:9092 --create --topic motivation --replication-factor 1
```

Overrides:

```yaml
image:
  repository: mihkels/kafka-tester-dev

imagePullSecrets:
  - regcred

bootstrapServers: "digital-ocean-cluster-kafka-bootstrap:9092"
topic: "motivation"

programmingLanguage:
  - "golang"
```

Debug container: 

```shell
docker run -it --entrypoint /bin/sh mihkels/kafka-tester-dev:producer-2024.05.09.13-golang
```