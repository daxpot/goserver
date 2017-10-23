package mypool

import (
	"time"
)

type Work struct {
	Fn        func(...interface{})
	Param 	  interface{}
	Completed bool
}

type Pool struct {
	queue chan int
	workers chan *Work
}

func New(size int) *Pool {
	if size <= 0 {
		size = 1
	}

	p := &Pool {
		queue: make(chan int, size),
		workers: make(chan *Work, size),
	}

	for i := 0; i<size; i++ {
		go func(p *Pool, i int) {
			for {
				t := <- p.workers
				t.Fn(t.Param, i)
				t.Completed = true
				<- p.queue
			}
		}(p, i)
	}
	return p
}

func (p *Pool) Add(fn func(...interface{}), param ...interface{}) {
	p.queue <- 1
	p.workers <- &Work{fn, param, false}
}

func (p *Pool) Wait() {
	for {
		if len(p.queue) == 0 {
			return
		}
		time.Sleep(time.Second)
	}
}

func (p *Pool) Length() int {
	return len(p.queue)
}