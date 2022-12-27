package main

import (
	"github.com/Shopify/sarama"
)

func setConfig(kafkaVersion sarama.KafkaVersion, cfg appConfig) *sarama.Config {
	config := sarama.NewConfig()
	config.Version = kafkaVersion

	if cfg.tls {
		config.Net.TLS.Enable = true
	}

	if cfg.user != "" && cfg.password != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = cfg.user
		config.Net.SASL.Password = cfg.password
	}

	if cfg.oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
	return config
}
