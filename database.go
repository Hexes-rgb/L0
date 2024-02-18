package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nats-io/stan.go"
)

func initializeTable() error {
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS orders (
		order_uid UUID PRIMARY KEY,
		order_data JSON NOT NULL
    )`)

	return err
}

func saveToDB(u uuid.UUID, msg *stan.Msg, ackChannel *chan *stan.Msg) {
	var conn *pgxpool.Conn
	var err error

	if conn, err = dbpool.Acquire(context.Background()); err != nil {
		log.Printf("Conn error: %v", err)
	}
	defer conn.Release()

	if _, err = conn.Exec(context.Background(), "INSERT INTO orders(order_uid, order_data) VALUES ($1, $2)", u, msg.Data); err != nil {
		log.Printf("Order %v saving failed: %v", u, err)
	} else {
		log.Printf("Order %v create successfully", u)
		*ackChannel <- msg
	}
}
