package database

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func Open() *pgx.Conn {

	var conn *pgx.Conn
	var err error

	if gin.Mode() == gin.ReleaseMode {

		conn, err = pgx.Connect(context.Background(), "postgres://postgres:q123a456@localhost:5432/api")
	} else {
		conn, err = pgx.Connect(context.Background(), "postgres://postgres:q123a456@localhost:5432/api")
	}

	if err != nil {
		panic(err)
	}

	return conn

}
