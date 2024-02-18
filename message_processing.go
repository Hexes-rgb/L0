package main

import (
	"log"

	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
)

func processMessages(msgChannel *chan *stan.Msg, ackChannel *chan *stan.Msg) {
	for msg := range *msgChannel {
		if validateData(msg) {
			u := uuid.New()
			// saveToCache(msg)
			go saveToDB(u, msg, ackChannel)
		} else {
			log.Println("Received invalid data")
			*ackChannel <- msg
		}
	}
}
func processAckQueue(ackChannel *chan *stan.Msg) {
	for msg := range *ackChannel {
		if err := msg.Ack(); err != nil {
			log.Printf("Message confirmation error: %v", err)
		}
	}
}
