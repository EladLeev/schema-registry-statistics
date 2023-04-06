package main

import (
	"context"
	"encoding/binary"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/eladleev/schema-registry-statistics/utils"
)

type Consumer struct {
	ready        chan bool
	stats        utils.ResultStats
	config       appConfig
	consumerLock sync.RWMutex
}

func main() {
	cfg := parseFlags()
	log.SetPrefix("[sr-stats] ")
	log.Printf("Starting to consume from %v", cfg.topic)

	version, err := sarama.ParseKafkaVersion(cfg.version)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	config := setConfig(version, cfg)

	consumer := Consumer{
		ready: make(chan bool),
		stats: utils.ResultStats{
			StatMap:     map[string]int{"TOTAL": 0},
			ResultStore: map[uint32][]int{},
		},
		config:       cfg,
		consumerLock: sync.RWMutex{},
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(strings.Split(cfg.bootstrapServers, ","), cfg.group, config)
	if err != nil {
		log.Panicf("Error creating consumer group: %v", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, strings.Split(cfg.topic, ","), &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready
	log.Println("Consumer up and running!...")
	log.Println("Use SIGINT to stop consuming.")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	case <-sigterm:
		log.Println("terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	consumedMessages := consumer.stats.StatMap["TOTAL"]
	log.Printf("Total messages consumed: %v\n", consumedMessages)
	for k, v := range consumer.stats.StatMap {
		if k == "TOTAL" {
			continue
		} else if k == "ERROR" {
			defer log.Printf("Unable to decode schema in %v messages. They might be empty, or do not contains any schema.", v)
		} else {
			utils.CalcPercentile(k, v, consumedMessages)
		}
	}
	if consumer.config.store {
		utils.DumpStats(consumer.stats, consumer.config.path)
	}
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if len(message.Value) < 5 {
				log.Printf("error encoding message offset: %v\n", message.Offset)
				// append error value, using 4294967295 as a dummy value (end of uint32)
				utils.CalcStat(consumer.stats, 4294967295, &consumer.consumerLock)
				break
			}
			schemaId := binary.BigEndian.Uint32(message.Value[1:5])
			utils.CalcStat(consumer.stats, schemaId, &consumer.consumerLock)
			if consumer.config.store { // lock map, and build result for analysis
				utils.AppendResult(consumer.stats, message.Offset, schemaId, &consumer.consumerLock)
			}
			consumer.consumerLock.RLock()
			if consumer.stats.StatMap["TOTAL"]%100 == 0 { // I'm still alive :)
				log.Printf("acked 100 messages\n")
			}
			consumer.consumerLock.RUnlock()
			// ack
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}
