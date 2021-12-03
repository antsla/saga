package broker

import (
	"context"
	"fmt"
	"os"

	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"
)

func RunConsumers(ctx context.Context, handlers map[string]sarama.ConsumerGroupHandler) {
	kafkaConsumerGroups := initAllConsumerGroups()

	for topic, group := range kafkaConsumerGroups {
		go func(topic string, group *sarama.ConsumerGroup) {
			defer func() {
				if r := recover(); r != nil {
					log.Error().Str("panic", "true").Msg(fmt.Sprintf("%s", r))
				}
			}()

			for {
				err := (*group).Consume(ctx, []string{topic}, handlers[topic])
				if err != nil {
					log.Error().Err(err).Msg("consumer group error")
				}
			}
		}(topic, group)
	}
}

func initAllConsumerGroups() map[string]*sarama.ConsumerGroup {
	return map[string]*sarama.ConsumerGroup{
		os.Getenv("GOODS_CREATED_TOPIC"):  initGroup(os.Getenv("GOODS_CREATED_TOPIC")),
		os.Getenv("GOODS_REJECTED_TOPIC"): initGroup(os.Getenv("GOODS_REJECTED_TOPIC")),
	}
}

func initGroup(topic string) *sarama.ConsumerGroup {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_3_0_0
	cfg.Consumer.Return.Errors = true

	group, err := sarama.NewConsumerGroup([]string{os.Getenv("KAFKA_ADDR")}, topic, cfg)
	if err != nil {
		log.Error().Err(err).Msg("Message hasn't been marshaled.")
		return nil
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error().Str("panic", "true").Msg(fmt.Sprintf("%s", r))
			}
		}()

		for err := range group.Errors() {
			log.Error().Err(err).Msg("consumer group error")
		}
	}()

	return &group
}
