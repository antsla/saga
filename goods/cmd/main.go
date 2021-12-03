package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Shopify/sarama"
	"github.com/antsla/goods/pkg/broker"
	"github.com/antsla/goods/pkg/datastore"
	"github.com/antsla/goods/transport"
	"github.com/rs/zerolog/log"
)

func main() {
	db := datastore.InitDB()
	producer := broker.InitKafkaProducer()

	handlers := map[string]sarama.ConsumerGroupHandler{
		os.Getenv("ORDER_CREATED_TOPIC"): broker.BuildOrderCreatedHandler(db, producer),
	}
	broker.RunConsumers(context.Background(), handlers)

	server := transport.NewServer()
	fmt.Println("server is starting...")
	err := server.Start()
	if err != nil {
		log.Error().Err(err).Msg("Server hasn't been started.")
		os.Exit(1)
	}
}
