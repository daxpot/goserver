package mypool

import (
	"time"
)

type Work struct {
	Fn        func(...interface{})
	Param 	  interface{}
	Completed bool
}

type pool struct {
	queue chan int
	workers chan *Work
}

func New(size int) *pool {
	if size <= 0 {
		size = 1
	}

	p := &pool {
		queue: make(chan int, size),
		workers: make(chan *Work, size),
	}

	for i := 0; i<size; i++ {
		go func(p *pool, i int) {
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

func (p *pool) Add(fn func(...interface{}), param ...interface{}) {
	p.queue <- 1
	p.workers <- &Work{fn, param, false}
}

func (p *pool) Wait() {
	for {
		if len(p.queue) == 0 {
			return
		}
		time.Sleep(time.Second)
	}
}