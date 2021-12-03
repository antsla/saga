package broker

import (
	"os"

	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"
)

func InitKafkaProducer() sarama.SyncProducer {
	brokerCfg := sarama.NewConfig()
	brokerCfg.Producer.RequiredAcks = sarama.WaitForAll
	brokerCfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{os.Getenv("KAFKA_ADDR")}, brokerCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Kafka error.")
		os.Exit(1)
	}

	return producer
}
