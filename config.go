package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/nats-io/stan.go"
)

type ConnPoolConfig struct {
	ConnStr         string
	MaxConns        int32
	MaxConnIdleTime int
}
type Config struct {
	StanConnOpts        []stan.Option
	StanSubOpts         []stan.SubscriptionOption
	StanClusterName     string
	StanClientID        string
	StanChannelName     string
	ValidationSchema    string
	DispacherBufferSize int
	ConnPoolConfig      ConnPoolConfig
	MaxCacheSize        int
	WorkersCount        int
	ServicePort         string
}

func InitializeConfig() (*Config, error) {
	stanConnStr, err := getEnv("STAN_URL")
	if err != nil {
		return nil, fmt.Errorf("missing STAN_URL: %v", err)
	}

	clusterName, err := getEnv("STAN_CLUSTER_NAME")
	if err != nil {
		return nil, fmt.Errorf("missing STAN_CLUSTER_NAME: %v", err)
	}

	servicePort, err := getEnv("SERVICE_PORT")
	if err != nil {
		return nil, fmt.Errorf("missing SERVICE_PORT: %v", err)
	}

	clientID, err := getEnv("STAN_CLIENT_ID")
	if err != nil {
		return nil, fmt.Errorf("missing STAN_CLIENT_ID: %v", err)
	}

	channelName, err := getEnv("STAN_CHANNEL_NAME")
	if err != nil {
		return nil, fmt.Errorf("missing STAN_CHANNEL_NAME: %v", err)
	}

	durableName, err := getEnv("STAN_DURABLE_NAME")
	if err != nil {
		return nil, fmt.Errorf("missing STAN_DURABLE_NAME: %v", err)
	}

	maxInflight, err := getUint64Env("STAN_MAX_INFLIGHT")
	if err != nil {
		return nil, fmt.Errorf("missing STAN_MAX_INFLIGHT: %v", err)
	}

	pgConn, err := getEnv("PG_CONN_STR")
	if err != nil {
		return nil, fmt.Errorf("missing PG_CONN_STR: %v", err)
	}

	maxPubAcksInflight, err := getUint64Env("STAN_MAX_PUB_ACKS_INFLIGHT")
	if err != nil {
		return nil, fmt.Errorf("missing STAN_MAX_PUB_ACKS_INFLIGHT: %v", err)
	}

	pubAckWait, err := getUint64Env("STAN_PUB_ACK_WAIT")
	if err != nil {
		return nil, fmt.Errorf("missing STAN_PUB_ACK_WAIT: %v", err)
	}

	dispatcherBufferSize, err := getUint64Env("DISPATCHER_BUFFER_SIZE")
	if err != nil {
		return nil, fmt.Errorf("missing DISPATCHER_BUFFER_SIZE: %v", err)
	}

	workersCount, err := getUint64Env("WORKERS_COUNT")
	if err != nil {
		return nil, fmt.Errorf("missing WORKERS_COUNT: %v", err)
	}

	maxConns, err := getUint64Env("PGX_POOL_MAX_CONNS")
	if err != nil {
		return nil, fmt.Errorf("missing PGX_POOL_MAX_CONNS: %v", err)
	}

	connIdleTime, err := getUint64Env("PGX_POOL_MAX_IDLE_TIME")
	if err != nil {
		return nil, fmt.Errorf("missing PGX_POOL_MAX_IDLE_TIME: %v", err)
	}

	maxCacheSize, err := getUint64Env("CACHE_MAX_SIZE")
	if err != nil {
		return nil, fmt.Errorf("missing CACHE_MAX_SIZE: %v", err)
	}

	opts := []stan.Option{
		stan.NatsURL(stanConnStr),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Printf("Connection lost, reason: %v", reason)
		}),
		stan.MaxPubAcksInflight(int(maxPubAcksInflight)),
		stan.PubAckWait(time.Duration(pubAckWait) * time.Second),
	}
	subOpts := []stan.SubscriptionOption{
		stan.DurableName(durableName),
		stan.MaxInflight(int(maxInflight)),
		stan.SetManualAckMode(),
		stan.DeliverAllAvailable(),
	}
	connPoolConfig := ConnPoolConfig{
		ConnStr:         pgConn,
		MaxConns:        int32(maxConns),
		MaxConnIdleTime: int(connIdleTime),
	}
	return &Config{
		StanConnOpts:        opts,
		StanSubOpts:         subOpts,
		StanClusterName:     clusterName,
		StanClientID:        clientID,
		StanChannelName:     channelName,
		ValidationSchema:    schema,
		DispacherBufferSize: int(dispatcherBufferSize),
		ConnPoolConfig:      connPoolConfig,
		MaxCacheSize:        int(maxCacheSize),
		ServicePort:         servicePort,
		WorkersCount:        int(workersCount),
	}, nil
}

func getEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("environment variable %s is empty", key)
	}
	return value, nil
}

func getUint64Env(key string) (uint64, error) {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return 0, fmt.Errorf("environment variable %s is empty", key)
	}
	value, err := strconv.ParseUint(valueStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse environment variable %s: %v", key, err)
	}
	return value, nil
}
