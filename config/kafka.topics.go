package config

import "education-flow/internal/kafka"

// Def.
var kafkaConf = kafka.Config{
	CaPath:   ``,
	CertPath: ``,
	KeyPath:  ``,
	Brokers:  `localhost:9092`,
}

// TopicFileStatusUpdate is the Kafka topic for file status updates.
const (
	TopicFileStatusUpdate = "storage.file.status.update"
)
