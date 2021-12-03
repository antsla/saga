package broker

import (
	"context"
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

type GoodsCreatedEvent struct {
	Data struct {
		OrderID int64 `json:"order_id"`
	} `json:"data"`
}

type GoodsCreatedHandler struct {
	db *pgxpool.Pool
}

func BuildGoodsCreatedHandler(db *pgxpool.Pool) GoodsCreatedHandler {
	return GoodsCreatedHandler{db: db}
}

func (gch GoodsCreatedHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		gce := GoodsCreatedEvent{}
		err := json.Unmarshal(msg.Value, &gce)
		if err != nil {
			log.Error().Err(err).Msg("Event hasn't been handled.")
			session.MarkMessage(msg, "")
			continue
		}

		_, err = gch.db.Exec(context.Background(), `UPDATE orders SET status_id = 2 WHERE id = $1`, gce.Data.OrderID)
		if err != nil {
			log.Error().Err(err).Msg("Event hasn't been inserted.")
		}

		session.MarkMessage(msg, "")
	}

	return nil
}

func (GoodsCreatedHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (GoodsCreatedHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
