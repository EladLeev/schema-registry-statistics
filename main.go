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

	"github.com/IBM/sarama"
	"github.com/eladleev/schema-registry-statistics/utils"
	"github.com/fatih/color"
)

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

// Setup is setting up a new CG session
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

// Cleanup function to clean after the CG
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	consumedMessages := consumer.stats.StatMap["TOTAL"]
	log.Printf("Total messages consumed: %v\n", consumedMessages)
	utils.BuildPercentileMap(consumer.stats.StatMap)

	// Print results
	for k, v := range utils.PercentileMap {
		c := color.New(color.FgGreen)
		c.Sprintf("Schema ID %v => %v%%\n", k, v)
	}

	// Dump stats to file
	if consumer.config.store {
		utils.DumpStats(consumer.stats, consumer.config.path)
	}

	// Build Charts
	if consumer.config.chart {
		utils.GenChart()
	}
	return nil
}

// ConsumeClaim processes Kafka messages from a given topic and partition within a consumer group.
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
			schemaID := binary.BigEndian.Uint32(message.Value[1:5])
			utils.CalcStat(consumer.stats, schemaID, &consumer.consumerLock)
			if consumer.config.store { // lock map, and build result for analysis
				utils.AppendResult(consumer.stats, message.Offset, schemaID, &consumer.consumerLock)
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
