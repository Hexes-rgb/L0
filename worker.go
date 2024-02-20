package main

import (
	"log"
)

type Worker struct {
	ID   int
	Quit chan bool
}

func NewWorker(id int) *Worker {
	return &Worker{
		ID:   id,
		Quit: make(chan bool),
	}
}

func (w *Worker) Start(taskChan <-chan interface{}) {
	go func() {
		for {
			select {
			case task := <-taskChan:
				switch t := task.(type) {
				case SaveToDBTask:
					go saveToDB(&t)
				case CacheTask:
					go cache.Add(t.data)
				case AckMsgTask:
					if err := t.msg.Ack(); err != nil {
						log.Printf("Message confirmation error: %v", err)
					}
				}

			case <-w.Quit:
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.Quit <- true
	}()
}

type Dispatcher struct {
	WorkerPools []*Worker
	TaskChan    chan interface{}
	MaxWorkers  int
	BufferSize  int
}

func NewDispatcher(maxWorkers, bufferSize int) *Dispatcher {
	return &Dispatcher{
		WorkerPools: make([]*Worker, maxWorkers),
		TaskChan:    make(chan interface{}, bufferSize),
		MaxWorkers:  maxWorkers,
		BufferSize:  bufferSize,
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(i + 1)
		worker.Start(d.TaskChan)
		d.WorkerPools[i] = worker
	}
}

func (d *Dispatcher) AddTask(task interface{}) {
	d.TaskChan <- task
}
