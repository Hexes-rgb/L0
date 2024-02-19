package main

import (
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
	StanConnOpts     []stan.Option
	StanSubOpts      []stan.SubscriptionOption
	StanClusterName  string
	StanClientID     string
	StanChannelName  string
	ValidationSchema string
	MsgChannelSize   uint64
	ConnPoolConfig   ConnPoolConfig
	MaxCacheSize     int
}

var schema = `
{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"type": "object",
	"properties": {
	  "order_uid": {"type": "string"},
	  "track_number": {"type": "string"},
	  "entry": {"type": "string"},
	  "delivery": {
		"type": "object",
		"properties": {
		  "name": {"type": "string"},
		  "phone": {"type": "string"},
		  "zip": {"type": "string"},
		  "city": {"type": "string"},
		  "address": {"type": "string"},
		  "region": {"type": "string"},
		  "email": {"type": "string", "format": "email"}
		},
		"required": ["name", "phone", "zip", "city", "address", "region", "email"]
	  },
	  "payment": {
		"type": "object",
		"properties": {
		  "transaction": {"type": "string"},
		  "request_id": {"type": "string"},
		  "currency": {"type": "string"},
		  "provider": {"type": "string"},
		  "amount": {"type": "integer"},
		  "payment_dt": {"type": "integer"},
		  "bank": {"type": "string"},
		  "delivery_cost": {"type": "integer"},
		  "goods_total": {"type": "integer"},
		  "custom_fee": {"type": "integer"}
		},
		"required": ["transaction", "currency", "provider", "amount", "payment_dt", "bank", "delivery_cost", "goods_total", "custom_fee"]
	  },
	  "items": {
		"type": "array",
		"items": {
		  "type": "object",
		  "properties": {
			"chrt_id": {"type": "integer"},
			"track_number": {"type": "string"},
			"price": {"type": "integer"},
			"rid": {"type": "string"},
			"name": {"type": "string"},
			"sale": {"type": "integer"},
			"size": {"type": "string"},
			"total_price": {"type": "integer"},
			"nm_id": {"type": "integer"},
			"brand": {"type": "string"},
			"status": {"type": "integer"}
		  },
		  "required": ["chrt_id", "track_number", "price", "rid", "name", "sale", "size", "total_price", "nm_id", "brand", "status"]
		}
	  },
	  "locale": {"type": "string"},
	  "internal_signature": {"type": "string"},
	  "customer_id": {"type": "string"},
	  "delivery_service": {"type": "string"},
	  "shardkey": {"type": "string"},
	  "sm_id": {"type": "integer"},
	  "date_created": {"type": "string", "format": "date-time"},
	  "oof_shard": {"type": "string"}
	},
	"required": ["order_uid", "track_number", "entry", "delivery", "payment", "items", "locale", "customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard"]
  }
  `

func InitializeConfig() *Config {
	var clusterName, clientID, channelName, durableName, pgConn string
	var maxInflight, maxPubAcksInflight, pubAckWait, msgChannelSize, maxConns, connIdleTime, maxCacheSize uint64
	var err error
	if clusterName = os.Getenv("STAN_CLUSTER_NAME"); clusterName == "" {
		panic("Entry valid cluster id")
	}
	if clientID = os.Getenv("STAN_CLIENT_ID"); clientID == "" {
		panic("Entry valid client id")
	}
	if channelName = os.Getenv("STAN_CHANNEL_NAME"); channelName == "" {
		panic("Entry valid channel name")
	}
	if durableName = os.Getenv("STAN_DURABLE_NAME"); durableName == "" {
		panic("Entry valid durable name")
	}
	if maxInflight, err = strconv.ParseUint(os.Getenv("STAN_MAX_INFLIGHT"), 10, 64); err != nil {
		panic("Entry valid max inflight")
	}
	if pgConn = os.Getenv("PG_CONN_STR"); pgConn == "" {
		panic("Entry valid postgres connection string")
	}
	if maxPubAcksInflight, err = strconv.ParseUint(os.Getenv("STAN_MAX_PUB_ACKS_INFLIGHT"), 10, 64); err != nil {
		panic("Entry valid max pub acks inflight")
	}
	if pubAckWait, err = strconv.ParseUint(os.Getenv("STAN_PUB_ACK_WAIT"), 10, 64); err != nil {
		panic("Entry valid max pub acks wait")
	}
	if msgChannelSize, err = strconv.ParseUint(os.Getenv("MSG_CHANNEL_SIZE"), 10, 64); err != nil {
		panic("Entry valid msg channel size")
	}
	if maxConns, err = strconv.ParseUint(os.Getenv("PGX_POOL_MAX_CONNS"), 10, 64); err != nil {
		panic("Entry valid max conns count")
	}
	if connIdleTime, err = strconv.ParseUint(os.Getenv("PGX_POOL_MAX_IDLE_TIME"), 10, 64); err != nil {
		panic("Entry valid conn idle time in seconds")
	}
	if maxCacheSize, err = strconv.ParseUint(os.Getenv("CACHE_MAX_SIZE"), 10, 64); err != nil {
		panic("Entry valid max cache size")
	}
	opts := []stan.Option{
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
		StanConnOpts:     opts,
		StanSubOpts:      subOpts,
		StanClusterName:  clusterName,
		StanClientID:     clientID,
		StanChannelName:  channelName,
		ValidationSchema: schema,
		MsgChannelSize:   msgChannelSize,
		ConnPoolConfig:   connPoolConfig,
		MaxCacheSize:     int(maxCacheSize),
	}

}
