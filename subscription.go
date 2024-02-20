package main

import (
	"log"

	"github.com/nats-io/stan.go"
)

func subscribeToSTAN() {
	handler := func(msg *stan.Msg) {
		if validateData(msg) {
			createSaveToDBTask(msg)
		} else {
			log.Println("Received invalid data")
			createAckMsgTask(msg)
		}
	}

	_, err := sc.Subscribe(config.StanChannelName, handler, config.StanSubOpts...)
	if err != nil {
		log.Fatalf("Channel subscription error: %v", err)
	}
}
