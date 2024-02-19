package main

import (
	"log"
)

func processSaveToDBQueue() {
	for msg := range msgChannel {
		if validateData(msg) {
			go saveToDB(msg)
		} else {
			log.Println("Received invalid data")
			ackChannel <- msg
		}
	}
}

func processAckQueue() {
	for msg := range ackChannel {
		if err := msg.Ack(); err != nil {
			log.Printf("Message confirmation error: %v", err)
		}
	}
}

func processDataCacheQueue() {
	for data := range dataToCacheChannel {
		go cache.Add(data)
	}
}
