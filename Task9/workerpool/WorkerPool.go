package workerpool

import (
	"context"
	"sync"
	"time"
)

// WPool - struct for getting work for goroutines, and wait finish work
type WPool struct {
	tasks chan func(context.Context)
	wait  sync.WaitGroup
}

// Workers - func create new forkers in struct for wait finish work, and start goroutine
func Workers(count int) *WPool {
	pool := &WPool{
		tasks: make(chan func(context.Context), 100),
	}
	for i := 0; i < count; i++ {
		pool.wait.Add(1)
		go pool.worker(i)
	}
	return pool
}

// Submit - give task for goroutine
func (p *WPool) Submit(task func(ctx context.Context)) {
	p.tasks <- task
}

// worker - start task and waiting for complete
func (p *WPool) worker(id int) {
	defer p.wait.Done()
	id += 0
	for task := range p.tasks {
		ctx := context.Background()
		task(ctx)
		time.Sleep(10 * time.Millisecond)
	}
}

// Shutdown - close channel
func (p *WPool) Shutdown() {
	close(p.tasks)
	p.wait.Wait()
}
