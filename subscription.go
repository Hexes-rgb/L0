package main

import (
	"log"

	"github.com/nats-io/stan.go"
)

func subscribeToSTAN(sc stan.Conn, msgChannel *chan *stan.Msg) {
	handler := func(msg *stan.Msg) {
		*msgChannel <- msg
	}

	_, err := sc.Subscribe(config.StanChannelName, handler, config.StanSubOpts...)
	if err != nil {
		log.Fatalf("Channel subscription error: %v", err)
	}
}
