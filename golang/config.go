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

	Group                  string
	EnableStatistics       bool
	StatisticsCollectorURL string
	WorkerName             string

	UseHeaders bool
}

var ConfigInstance *Config

func init() {
	ConfigInstance = NewConfig()
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Could not load .env file will fallback to environment variables")
	}

	return &Config{
		ApplicationMode:        getEnvString("APPLICATION_MODE", "producer"),
		Brokers:                getEnvString("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092"),
		Topic:                  getEnvString("TOPIC", "test"),
		ProducerDataFile:       getEnvString("MOTIVATION_FILE", "data.csv"),
		NumberOfSamples:        getEnvInt("NUM_SAMPLES", 100),
		BatchSize:              getEnvInt("BATCH_SIZE", 10),
		SleepInterval:          time.Duration(getEnvInt("SLEEP_TIME", 10)),
		Group:                  getEnvString("CONSUMER_GROUP", "golang-test-cg"),
		EnableStatistics:       getEnvString("ENABLE_STATISTICS", "false") == "true",
		StatisticsCollectorURL: getEnvString("STATISTICS_COLLECTOR_URL", "http://localhost:8080"),
		UseHeaders:             getEnvBool("ENABLE_HEADERS", false),
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

func getEnvBool(key string, defaultValue bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	parsedValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("Warning: Could not parse boolean value for key %s, using default %v", key, defaultValue)
		return defaultValue
	}

	return parsedValue
}
