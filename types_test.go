package main

import (
	"sync"
	"testing"
	"time"

	"github.com/eladleev/schema-registry-statistics/utils"
	"gotest.tools/assert"
)

func TestConsumerInitialization(t *testing.T) {
	config := appConfig{
		bootstrapServers: "localhost:9092",
		version:          "2.6.0",
		group:            "my-group",
		topic:            "my-topic",
		user:             "my-user",
		password:         "my-password",
		tls:              true,
		caCert:           "ca.pem",
		path:             "/tmp",
		limit:            1000,
		oldest:           true,
		verbose:          true,
		store:            true,
		chart:            true,
	}

	consumer := Consumer{
		ready:        make(chan bool),
		stats:        utils.ResultStats{},
		config:       config,
		consumerLock: sync.RWMutex{},
	}

	assert.Equal(t, config.bootstrapServers, consumer.config.bootstrapServers)
	assert.Equal(t, config.version, consumer.config.version)
	assert.Equal(t, config.group, consumer.config.group)
	assert.Equal(t, config.topic, consumer.config.topic)
	assert.Equal(t, config.user, consumer.config.user)
	assert.Equal(t, config.password, consumer.config.password)
	assert.Equal(t, config.tls, consumer.config.tls)
	assert.Equal(t, config.caCert, consumer.config.caCert)
	assert.Equal(t, config.path, consumer.config.path)
	assert.Equal(t, config.limit, consumer.config.limit)
	assert.Equal(t, config.oldest, consumer.config.oldest)
	assert.Equal(t, config.verbose, consumer.config.verbose)
	assert.Equal(t, config.store, consumer.config.store)
	assert.Equal(t, config.chart, consumer.config.chart)
}

func TestConsumerLocking(t *testing.T) {
	config := appConfig{
		bootstrapServers: "localhost:9092",
		version:          "2.6.0",
		group:            "my-group",
		topic:            "my-topic",
		user:             "my-user",
		password:         "my-password",
		tls:              true,
		caCert:           "ca.pem",
		path:             "/tmp",
		limit:            1000,
		oldest:           true,
		verbose:          true,
		store:            true,
		chart:            true,
	}

	consumer := Consumer{
		ready:        make(chan bool),
		stats:        utils.ResultStats{},
		config:       config,
		consumerLock: sync.RWMutex{},
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		consumer.consumerLock.Lock()
		defer consumer.consumerLock.Unlock()

		time.Sleep(time.Millisecond * 100)
		wg.Done()
	}()

	go func() {
		consumer.consumerLock.Lock()
		defer consumer.consumerLock.Unlock()

		wg.Done()
	}()

	wg.Wait()
}
