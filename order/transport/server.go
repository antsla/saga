package transport

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Shopify/sarama"
	"github.com/antsla/order/pkg/model"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

type Server struct {
	router        *mux.Router
	db            *pgxpool.Pool
	kafkaProducer sarama.SyncProducer
}

func NewServer(db *pgxpool.Pool, kafkaProducer sarama.SyncProducer) Server {
	s := Server{}
	s.kafkaProducer = kafkaProducer
	s.db = db
	s.router = mux.NewRouter()

	s.router.HandleFunc("/v1/orders", s.CreateOrderV1).Methods(http.MethodPost)

	return s
}

func (s Server) Start() error {
	return http.ListenAndServe(":"+os.Getenv("HTTP_BIND"), s.router)
}

func (s Server) CreateOrderV1(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Data hasn't been read.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orderData := model.OrderData{}
	err = json.Unmarshal(body, &orderData)
	if err != nil {
		log.Error().Err(err).Msg("Data hasn't been parsed.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var orderID int64
	err = s.db.QueryRow(context.Background(), `INSERT INTO orders (user_id, status_id, created_at) VALUES ($1, 1, NOW()) RETURNING id`, orderData.UserID).Scan(&orderID)
	if err != nil {
		log.Error().Err(err).Msg("Order hasn't been created.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	msg := model.CreatedOrderMsg{Data: model.Order{
		ID:       orderID,
		GoodsIds: orderData.GoodsIds,
	}}
	msgStr, err := json.Marshal(msg)
	if err != nil {
		log.Error().Err(err).Msg("Message hasn't been marshaled.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	producerMsg := &sarama.ProducerMessage{Topic: os.Getenv("ORDER_CREATED_TOPIC"), Value: sarama.StringEncoder(msgStr)}
	_, _, err = s.kafkaProducer.SendMessage(producerMsg)
	if err != nil {
		log.Error().Err(err).Msg("Message hasn't been sent.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
