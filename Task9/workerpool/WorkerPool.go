package workerpool

import (
	"sync"
	"time"
)

// WPool - struct for getting work for goroutines, and wait finish work
type WPool struct {
	tasks chan func()
	wait  sync.WaitGroup
}

// Workers - func create new forkers in struct for wait finish work, and start goroutine
func Workers(count int) *WPool {
	pool := &WPool{
		tasks: make(chan func(), 100),
	}
	for i := 0; i < count; i++ {
		pool.wait.Add(1)
		go pool.worker(i)
	}
	return pool
}

// Submit - give task for goroutine
func (p *WPool) Submit(task func()) {
	p.tasks <- task
}

// worker - start task and waiting for complete
func (p *WPool) worker(id int) {
	defer p.wait.Done()
	for task := range p.tasks {
		task()
		time.Sleep(10 * time.Millisecond)
	}
}

// Shutdown - close channel
func (p *WPool) Shutdown() {
	close(p.tasks)
	p.wait.Wait()
}
