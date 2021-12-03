package broker

import (
	"context"
	"encoding/json"
	"os"

	"github.com/Shopify/sarama"
	"github.com/antsla/goods/pkg/model"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

type OrderCreatedEvent struct {
	Data struct {
		ID       int64   `json:"id"`
		GoodsIds []int64 `json:"goods_ids"`
	} `json:"data"`
}

type OrderCreatedHandler struct {
	db       *pgxpool.Pool
	producer sarama.SyncProducer
}

func BuildOrderCreatedHandler(db *pgxpool.Pool, producer sarama.SyncProducer) OrderCreatedHandler {
	return OrderCreatedHandler{db, producer}
}

func (och OrderCreatedHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		oce := OrderCreatedEvent{}
		err := json.Unmarshal(msg.Value, &oce)
		if err != nil {
			log.Error().Err(err).Msg("Event hasn't been handled.")
			session.MarkMessage(msg, "")
			continue
		}

		ctx := context.Background()
		tx, err := och.db.Begin(ctx)
		for _, goodsID := range oce.Data.GoodsIds {
			_, err = och.db.Exec(context.Background(), `INSERT INTO goods (goods_id, order_id, created_at) VALUES ($1, $2, NOW())`, goodsID, oce.Data.ID)
			if err != nil {
				tx.Rollback(ctx)
				log.Error().Err(err).Msg("Inserting error.")

				err := och.sendRejected(oce.Data.ID)
				if err != nil {
					log.Error().Err(err).Msg("Event hasn't been sent.")
				}
				session.MarkMessage(msg, "")
				continue
			}
		}

		err = tx.Commit(ctx)
		if err != nil {
			err := tx.Rollback(ctx)
			log.Error().Err(err).Msg("Transaction commit error.")
			rErr := och.sendRejected(oce.Data.ID)
			if rErr != nil {
				log.Error().Err(rErr).Msg("Event hasn't been sent.")
			}
			session.MarkMessage(msg, "")
			continue
		}

		err = och.sendCreated(oce.Data.ID)
		if err != nil {
			log.Error().Err(err).Msg("Event hasn't been sent.")
		}
		session.MarkMessage(msg, "")
	}

	return nil
}

func (och OrderCreatedHandler) sendRejected(orderID int64) error {
	msg := model.RejectedGoodsMsg{Data: model.Goods{
		OrderID: orderID,
	}}
	msgStr, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	producerMsg := &sarama.ProducerMessage{Topic: os.Getenv("GOODS_REJECTED_TOPIC"), Value: sarama.StringEncoder(msgStr)}
	_, _, err = och.producer.SendMessage(producerMsg)
	return err
}

func (och OrderCreatedHandler) sendCreated(orderID int64) error {
	msg := model.CreatedGoodsMsg{Data: model.Goods{
		OrderID: orderID,
	}}
	msgStr, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	producerMsg := &sarama.ProducerMessage{Topic: os.Getenv("GOODS_CREATED_TOPIC"), Value: sarama.StringEncoder(msgStr)}
	_, _, err = och.producer.SendMessage(producerMsg)
	return err
}

func (OrderCreatedHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (OrderCreatedHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
