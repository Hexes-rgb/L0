package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/nats-io/stan.go"
)

var config *Config
var dbpool *pgxpool.Pool

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
		panic(err)
	}

	config = InitializeConfig()
	//TODO: Pool config
	connConfig, err := pgxpool.ParseConfig(config.PGConnStr)
	if err != nil {
		log.Fatal("Error parsing connection string")
	}

	dbpool, err = pgxpool.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		log.Fatal("Error creating connection pool")
	}
	defer dbpool.Close()

	err = initializeTable()
	if err != nil {
		panic("DB is not initialized")
	}

	msgChannel := make(chan *stan.Msg, config.MsgChannelSize)
	ackChannel := make(chan *stan.Msg, config.MsgChannelSize)

	sc, err := stan.Connect(config.StanClusterName, config.StanClientID, config.StanConnOpts...)
	if err != nil {
		log.Fatalf("NATS Streaming connection error: %v", err)
		panic(err)
	}
	defer sc.Close()

	go subscribeToSTAN(sc, &msgChannel)
	go processMessages(&msgChannel, &ackChannel)
	go processAckQueue(&ackChannel)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}
