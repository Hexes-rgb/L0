package main

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func createConnsPool() error {
	connConfig, err := pgxpool.ParseConfig(config.ConnPoolConfig.ConnStr)
	if err != nil {
		return err
	}
	connConfig.MaxConns = config.ConnPoolConfig.MaxConns
	connConfig.MaxConnIdleTime = time.Duration(config.ConnPoolConfig.MaxConnIdleTime) * time.Second
	dbpool, err = pgxpool.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		return err
	}
	return nil
}

func initializeTable() error {
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		order_uid UUID UNIQUE,
		order_data JSONB NOT NULL,
		created_at TIMESTAMP NOT NULL
    )`)

	return err
}

func saveToDB(task *SaveToDBTask) {
	var conn *pgxpool.Conn
	var err error

	if conn, err = dbpool.Acquire(context.Background()); err != nil {
		log.Printf("Conn error: %v", err)
	}
	defer conn.Release()

	orderUID, jsonData := getUUIDFromJson(task.msg.Data)

	layout := "2006-01-02 15:04:05.000"
	timestamp := time.Now()

	timestampStr := timestamp.Format(layout)

	if _, err = conn.Exec(context.Background(), "INSERT INTO orders(order_uid, order_data, created_at) VALUES ($1, $2, $3)", orderUID, jsonData, timestampStr); err != nil {
		log.Printf("Order %v saving failed: %v", orderUID, err)
	} else {
		log.Printf("Order %v create successfully", orderUID)
		createAckMsgTask(task.msg)
		createCacheTask(NewDataToCache(orderUID, jsonData, timestamp))
	}
}

func getOrder(id string) ([]byte, bool) {
	var orderUID string
	var orderData []byte
	var createdAt time.Time

	err := dbpool.QueryRow(context.Background(), "SELECT order_uid, order_data, created_at FROM orders WHERE order_uid=$1", id).Scan(&orderUID, &orderData, &createdAt)
	if err != nil {
		return nil, false
	}

	data := &DataToCache{
		OrderUID:  orderUID,
		OrderData: orderData,
		CreatedAt: createdAt,
	}
	createCacheTask(data)
	return data.OrderData, true
}
