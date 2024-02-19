package main

import (
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
var cache *Cache
var dataToCacheChannel chan *DataToCache
var msgChannel, ackChannel chan *stan.Msg

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
		panic(err)
	}

	config = InitializeConfig()

	err := createConnsPool()
	if err != nil {
		log.Fatal("Error creating connection pool")
	}
	defer dbpool.Close()

	err = initializeTable()
	if err != nil {
		panic("DB is not initialized")
	}

	cache = NewCache(int(config.MaxCacheSize))

	if err = cache.InitializeCache(); err != nil {
		panic("Cache is not initialized")
	}

	msgChannel = make(chan *stan.Msg, config.MsgChannelSize)
	ackChannel = make(chan *stan.Msg, config.MsgChannelSize)
	dataToCacheChannel = make(chan *DataToCache, config.MsgChannelSize)

	sc, err := stan.Connect(config.StanClusterName, config.StanClientID, config.StanConnOpts...)
	if err != nil {
		log.Fatalf("NATS Streaming connection error: %v", err)
		panic(err)
	}
	defer sc.Close()

	go subscribeToSTAN(sc)
	go processSaveToDBQueue()
	go processDataCacheQueue()
	go createServer()
	go processAckQueue()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}
