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
					log.Printf("worker %d start task %T", w.ID, t)
				case CacheTask:
					go cache.Add(t.data)
					log.Printf("worker %d start task %T", w.ID, t)
				case AckMsgTask:
					log.Printf("worker %d start task %T", w.ID, t)
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
	TaskChan   chan interface{}
	MaxWorkers int
	BufferSize int
}

func NewDispatcher(maxWorkers, bufferSize int) *Dispatcher {
	return &Dispatcher{
		TaskChan:   make(chan interface{}, bufferSize),
		MaxWorkers: maxWorkers,
		BufferSize: bufferSize,
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(i + 1)
		worker.Start(d.TaskChan)
	}
}

func (d *Dispatcher) AddTask(task interface{}) {
	d.TaskChan <- task
}
