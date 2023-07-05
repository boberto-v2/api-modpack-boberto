package database

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4"
)

func OpenConnection() (*pgx.Conn, error, context.Context) {
	ctx := context.TODO()
	conn, err := pgx.Connect(ctx, os.Getenv("PG_URI"))
	if err != nil {
		panic(err)
	}
	err = conn.Ping(ctx)
	return conn, err, ctx
}
