package main

import (
	"sync"

	"github.com/eladleev/schema-registry-statistics/utils"
)

type appConfig struct {
	bootstrapServers, version, group, topic string
	user, password                          string
	tls                                     bool
	caCert                                  string
	path                                    string
	limit                                   int
	oldest                                  bool
	verbose                                 bool
	store                                   bool
	chart                                   bool
}

// Consumer represent a Kafka Consumer
type Consumer struct {
	ready        chan bool
	stats        utils.ResultStats
	config       appConfig
	consumerLock sync.RWMutex
}
