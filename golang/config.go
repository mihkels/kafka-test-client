package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ApplicationMode string

	Brokers string
	Topic   string

	ProducerDataFile string
	NumberOfSamples  int
	BatchSize        int
	SleepInterval    time.Duration

	Group string
}

var ConfigInstance *Config

func init() {
	ConfigInstance = NewConfig()
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		ApplicationMode:  getEnvString("APPLICATION_MODE", "producer"),
		Brokers:          getEnvString("BROKERS", "localhost:9092"),
		Topic:            getEnvString("TOPIC", "test"),
		ProducerDataFile: getEnvString("MOTIVATION_FILE", "data.csv"),
		NumberOfSamples:  getEnvInt("NUM_SAMPLES", 100),
		BatchSize:        getEnvInt("BATCH_SIZE", 10),
		SleepInterval:    time.Duration(getEnvInt("SLEEP_TIME", 10)),
		Group:            getEnvString("CONSUMER_GROUP", "golang-test-cg"),
	}
}

func getEnvString(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
