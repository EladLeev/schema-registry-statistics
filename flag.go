package main

import (
	"flag"
	"log"
	"os"

	"github.com/Shopify/sarama"
)

type appConfig struct {
	bootstrapServers, version, group, topic string
	user, password                          string
	tls                                     bool
	path                                    string
	limit                                   int
	oldest                                  bool
	verbose                                 bool
	store                                   bool
}

func parseFlags() appConfig {
	cfg := appConfig{}

	// Kafka configuration
	flag.StringVar(&cfg.bootstrapServers, "bootstrap", "localhost:9092", "The Kafka bootstrap servers, as a comma separated list")
	flag.StringVar(&cfg.group, "group", "schema-stats", "The Kafka consumer group name")
	flag.StringVar(&cfg.version, "version", "2.1.1", "The Kafka client version")
	flag.StringVar(&cfg.topic, "topic", "", "The Kafka topic to consume from")
	flag.StringVar(&cfg.user, "user", "", "The Kafka username")
	flag.StringVar(&cfg.password, "password", "", "The Kafka password")
	flag.BoolVar(&cfg.oldest, "oldest", true, "Consume from oldest offset")
	flag.BoolVar(&cfg.verbose, "verbose", false, "Switch to verbose logging")
	flag.IntVar(&cfg.limit, "limit", 0, "Limit consumer to N messages")
	flag.BoolVar(&cfg.tls, "tls", false, "Enable TLS connection")

	// Tool configuration
	flag.BoolVar(&cfg.store, "store", false, "Store results to file for analysis")
	flag.StringVar(&cfg.path, "path", "/tmp/results.json", "Default file to store the results")

	flag.Parse()

	if len(cfg.topic) == 0 {
		log.Fatal("No topic name was given. Please set the --topic flag and try again")
	}

	if cfg.verbose {
		sarama.Logger = log.New(os.Stdout, "[sr-stats DEBUG] ", log.LstdFlags)
	}

	return cfg
}
