package systems

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"

	_ "github.com/lib/pq"
)

func MustConn(cfg *AppConfig) *sqlx.DB {
	conn, err := sqlx.Open("postgres", cfg.Database.URL)

	if err != nil {
		log.Fatal().Err(err).Send()
	}

	return conn
}
