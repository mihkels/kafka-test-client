package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/IBM/sarama"
	"kafka_demo/internal/config"
	"log"
	"os"
)

type KafkaClient struct {
	Config *config.Config
}

func (k *KafkaClient) NewKafkaClientConfig() (*sarama.Config, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Version = sarama.V4_0_0_0

	if k.Config.UseTls {
		log.Println("Loading Kafka TLS configuration")
		if err := k.loadKafkaTLSConfig(kafkaConfig); err != nil {
			log.Println(err, "Failed to load Kafka TLS configuration")
			return nil, err
		}
		log.Println("Kafka TLS configuration loaded successfully")
	}

	return kafkaConfig, nil
}

func (k *KafkaClient) loadKafkaTLSConfig(kafkaConfig *sarama.Config) error {
	read := func(path string, label string) ([]byte, error) {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", label, err)
		}
		return b, nil
	}

	caCert, err := read(k.Config.CaCert, "CA cert")
	if err != nil {
		return err
	}
	userCert, err := read(k.Config.UserCert, "user cert")
	if err != nil {
		return err
	}
	userKey, err := read(k.Config.UserKey, "user key")
	if err != nil {
		return err
	}

	if len(caCert) == 0 || len(userCert) == 0 || len(userKey) == 0 {
		return fmt.Errorf("kafka user secret is missing required files")
	}

	cert, err := tls.X509KeyPair(userCert, userKey)
	if err != nil {
		return fmt.Errorf("failed to parse TLS certificate and key from Kafka user: %w", err)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCert) {
		return fmt.Errorf("failed to append CA certificate from Kafka user secret")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caPool,
		MinVersion:   tls.VersionTLS12, // Ensure TLS 1.2 or higher
	}

	if k.Config.AllowInsecure {
		tlsConfig.InsecureSkipVerify = true
	}

	kafkaConfig.Net.TLS.Enable = true
	kafkaConfig.Net.TLS.Config = tlsConfig

	return nil
}
