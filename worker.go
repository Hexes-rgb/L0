package main

import (
	"fmt"
)

type Worker struct {
	ID         int
	Task       chan Task
	WorkerPool chan chan Task
	Quit       chan bool
}

type Task struct {
	ID int
}

func NewWorker(id int, workerPool chan chan Task) Worker {
	return Worker{
		ID:         id,
		Task:       make(chan Task),
		WorkerPool: workerPool,
		Quit:       make(chan bool),
	}
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.Task

			select {
			case task := <-w.Task:
				fmt.Printf("Worker %d: started task %d\n", w.ID, task.ID)
				fmt.Printf("Worker %d: finished task %d\n", w.ID, task.ID)
			case <-w.Quit:
				fmt.Printf("Worker %d: stopping\n", w.ID)
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.Quit <- true
	}()
}

type Dispatcher struct {
	WorkerPool chan chan Task
	TaskQueue  chan Task
	MaxWorkers int
}

func NewDispatcher(maxWorkers, maxQueueSize int) *Dispatcher {
	pool := make(chan chan Task, maxWorkers)
	queue := make(chan Task, maxQueueSize)

	return &Dispatcher{
		WorkerPool: pool,
		TaskQueue:  queue,
		MaxWorkers: maxWorkers,
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(i+1, d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case task := <-d.TaskQueue:
			go func() {
				worker := <-d.WorkerPool
				worker <- task
			}()
		}
	}
}
