package main

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"path/filepath"

	"github.com/Shopify/sarama"
)

func loadKey(caFile string) *tls.Config {
	caCert, err := os.ReadFile(filepath.Clean(caFile))
	if err != nil {
		panic(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		RootCAs:    caCertPool,
	}
	return tlsConfig
}

func setConfig(kafkaVersion sarama.KafkaVersion, cfg appConfig) *sarama.Config {
	config := sarama.NewConfig()
	config.Version = kafkaVersion

	if cfg.tls {
		config.Net.TLS.Enable = true
		tlsCfg := loadKey(cfg.caCert)
		config.Net.TLS.Config = tlsCfg
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
