package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

var DBQueries *Queries
var Pool *pgxpool.Pool

func InitDb(context context.Context) error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		fmt.Errorf("DATABASE_URL is not set")
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		fmt.Println(err.Error())
	}

	pool, err := pgxpool.NewWithConfig(context, config)
	if err != nil {
		fmt.Println(err.Error())
	}

	queries := New(pool)
	DBQueries = queries
	Pool = pool

	return nil
}
