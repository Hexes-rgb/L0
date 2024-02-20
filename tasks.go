package main

import (
	"github.com/nats-io/stan.go"
)

type TaskType int

const (
	CacheTaskType TaskType = iota
	SaveToDBTaskType
	AckMsgTaskType
)

type SaveToDBTask struct {
	Type TaskType
	msg  *stan.Msg
}

type AckMsgTask struct {
	Type TaskType
	msg  *stan.Msg
}

type CacheTask struct {
	Type TaskType
	data *DataToCache
}

func createSaveToDBTask(msg *stan.Msg) {
	t := SaveToDBTask{
		Type: SaveToDBTaskType,
		msg:  msg,
	}
	dispatcher.AddTask(t)
}

func createCacheTask(data *DataToCache) {
	t := CacheTask{
		Type: CacheTaskType,
		data: data,
	}
	dispatcher.AddTask(t)
}

func createAckMsgTask(msg *stan.Msg) {
	t := AckMsgTask{
		Type: SaveToDBTaskType,
		msg:  msg,
	}
	dispatcher.AddTask(t)
}
