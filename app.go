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
var sc stan.Conn
var dispatcher *Dispatcher

func main() {
	var err error
	if err = godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if config, err = InitializeConfig(); err != nil {
		log.Fatalf("Error initializing config: %v", err)
	}

	if err = createConnsPool(); err != nil {
		log.Fatal("Error creating connection pool")
	}
	defer dbpool.Close()

	if err = initializeTable(); err != nil {
		log.Fatal("DB is not initialized")
	}

	cache = NewCache(config.MaxCacheSize)

	if err = cache.InitializeCache(); err != nil {
		log.Fatal("Cache is not initialized")
	}

	if sc, err = stan.Connect(config.StanClusterName, config.StanClientID, config.StanConnOpts...); err != nil {
		log.Fatalf("NATS Streaming connection error: %v", err)
	}
	defer sc.Close()

	dispatcher = NewDispatcher(config.WorkersCount, config.DispacherBufferSize)

	dispatcher.Run()

	go subscribeToSTAN()

	go createServer()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}
