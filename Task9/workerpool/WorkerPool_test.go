package workerpool

import (
	"context"
	"reflect"
	"sync"
	"testing"
)

func TestWPool_Shutdown(t *testing.T) {
	type fields struct {
		tasks chan func(context.Context)
		wait  sync.WaitGroup
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Shutdown tasks",
			fields: fields{
				tasks: make(chan func(context.Context), 1),
				wait:  sync.WaitGroup{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &WPool{
				tasks: tt.fields.tasks,
				wait:  tt.fields.wait,
			}

			go func() {
				p.tasks <- func(context.Context) {}
			}()

			p.Shutdown()
			select {
			case _, ok := <-p.tasks:
				if ok {
					t.Errorf("Channel wasnt closed")
				}
			default:
				t.Errorf("No task processed, channel may be blocked")
			}
		})
	}
}

func TestWPool_worker(t *testing.T) {
	type fields struct {
		tasks chan func(context.Context)
		wait  sync.WaitGroup
	}
	type args struct {
		id int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Worker execution",
			fields: fields{
				tasks: make(chan func(context.Context), 1),
				wait:  sync.WaitGroup{},
			},
			args: args{id: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &WPool{
				tasks: tt.fields.tasks,
				wait:  tt.fields.wait,
			}
			p.wait.Add(1)
			go p.worker(tt.args.id)
			p.tasks <- func(ctx context.Context) { t.Log("Worker executed task") }
			p.Shutdown()
		})
	}
}

func TestWorkers(t *testing.T) {
	type args struct {
		count int
	}
	tests := []struct {
		name string
		args args
		want *WPool
	}{
		{
			name: "Create pool with 3 workers",
			args: args{
				count: 3,
			},
			want: &WPool{
				tasks: make(chan func(context.Context), 100),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Workers(tt.args.count)
			if cap(got.tasks) != 100 {
				t.Errorf("Expected task channel capacity to be 100, got %v", cap(got.tasks))
			}
			if reflect.TypeOf(got.tasks).Kind() != reflect.Chan {
				t.Errorf("Expected tasks to be of type channel")
			}
		})
	}
}
