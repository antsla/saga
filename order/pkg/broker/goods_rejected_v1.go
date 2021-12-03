package broker

import (
	"context"
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

type GoodsRejectedEvent struct {
	Data struct {
		OrderID int64 `json:"order_id"`
	} `json:"data"`
}

type GoodsRejectedHandler struct {
	db *pgxpool.Pool
}

func BuildGoodsRejectedHandler(db *pgxpool.Pool) GoodsRejectedHandler {
	return GoodsRejectedHandler{db: db}
}

func (grh GoodsRejectedHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		gre := GoodsRejectedEvent{}
		err := json.Unmarshal(msg.Value, &gre)
		if err != nil {
			log.Error().Err(err).Msg("Event hasn't been handled.")
			session.MarkMessage(msg, "")
			continue
		}

		_, err = grh.db.Exec(context.Background(), `DELETE FROM orders WHERE id = $1`, gre.Data.OrderID)
		if err != nil {
			log.Error().Err(err).Msg("Event hasn't been inserted.")
		}

		session.MarkMessage(msg, "")
	}

	return nil
}

func (GoodsRejectedHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (GoodsRejectedHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
