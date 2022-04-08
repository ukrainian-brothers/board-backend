package common

import (
	"database/sql"
	"fmt"
	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq"
)

func InitPostgres(cfg *PostgresConfig) (*gorp.DbMap, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("cannot open database: %w", err)
	}
	return &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}, nil
}
