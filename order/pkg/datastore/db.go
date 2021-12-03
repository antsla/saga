package datastore

import (
	"context"
	"os"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

func InitDB() *pgxpool.Pool {
	cfg, _ := pgxpool.ParseConfig("")
	cfg.ConnConfig.Host = os.Getenv("HOST_DB")
	portInt, err := strconv.ParseInt(os.Getenv("PORT_DB"), 0, 16)
	cfg.ConnConfig.Port = uint16(portInt)
	cfg.ConnConfig.User = os.Getenv("POSTGRES_USER")
	cfg.ConnConfig.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.ConnConfig.Database = os.Getenv("POSTGRES_DB")
	cfg.ConnConfig.PreferSimpleProtocol = true
	cfg.MaxConns = 20

	dbPool, err := pgxpool.ConnectConfig(context.Background(), cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("DB error")
		os.Exit(1)
	}

	return dbPool
}
