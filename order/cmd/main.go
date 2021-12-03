package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Shopify/sarama"
	"github.com/antsla/order/pkg/broker"
	"github.com/antsla/order/pkg/datastore"
	"github.com/antsla/order/transport"
	"github.com/rs/zerolog/log"
)

func main() {
	db := datastore.InitDB()
	producer := broker.InitKafkaProducer()

	handlers := map[string]sarama.ConsumerGroupHandler{
		os.Getenv("GOODS_CREATED_TOPIC"):  broker.BuildGoodsCreatedHandler(db),
		os.Getenv("GOODS_REJECTED_TOPIC"): broker.BuildGoodsRejectedHandler(db),
	}
	broker.RunConsumers(context.Background(), handlers)

	server := transport.NewServer(db, producer)
	fmt.Println("server is starting...")
	err := server.Start()
	if err != nil {
		log.Error().Err(err).Msg("Server hasn't been started.")
		os.Exit(1)
	}
}
