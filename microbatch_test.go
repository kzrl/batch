package microbatch

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

var batchSize = 10
var maxAge = time.Second * 10

func TestMicrobatch_Submit(t *testing.T) {

	testJob := NewJob(42)
	var emptyJob Job
	var numWorkers = 2	
	tests := []struct {
		name   string
		job Job
		want   JobResult
	}{
		{name: "Submit a good job", job: testJob, want: JobResult{Job: testJob, Output: testJob.String()}},
		{name: "Submit an empty job", job: emptyJob, want: JobResult{Job: emptyJob, Output: emptyJob.String()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var nbp NoopBatchProcessor
			mb:= New(batchSize, maxAge, numWorkers, &nbp)
		
			if got := mb.Submit(tt.job); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Microbatch.Submit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMicrobatch_Start(t *testing.T) {
	tests := []struct {
		name   string
		want int
		numWorkers int
	}{
		{name: "Test Start() with 2 workers", numWorkers:2, want: 2},
		{name: "Test Start() with -1 workers", numWorkers: -1, want: 2},
		{name: "Test Start() with 20 workers", numWorkers: 20, want: 20},
		{name: "Test Start() with 0 workers", numWorkers: 0, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var nbp NoopBatchProcessor
			mb:= New(batchSize, maxAge, tt.numWorkers, &nbp)
			got := mb.Start()
			if got != tt.numWorkers {
				t.Errorf("Microbatch.Start() = %v, want %v", got, tt.want)	
			}
			mb.Shutdown()
		})
	}
}

// Not a very exciting test, but checks that it compiles OK and doesn't panic.
func TestMicrobatch_SubmitWorker(t *testing.T) {
	var numWorkers = 2	
	tests := []struct {
		name   string
		id int
	}{
		{name:"Create a worker with a positive ID", id: 0},
		{name:"Create a worker with a negative ID", id: -42},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var nbp NoopBatchProcessor
			mb:= New(batchSize, maxAge, numWorkers, &nbp)
			mb.Start()
			mb.Shutdown()
		})
	}
}

func TestMicrobatch_IsShutdown(t *testing.T) {
	var numWorkers = 2	
	tests := []struct {
		name   string
		start bool
		shutdown bool
		want   bool
	}{
		{name: "Start it up and shut it down again", start: true, shutdown: true, want: true},
		{name: "Start it up and leave it running", start: true, shutdown: false, want: false},	
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var nbp NoopBatchProcessor
			mb:= New(batchSize, maxAge, numWorkers, &nbp)
			if tt.start {
				mb.Start()
			}
			if tt.shutdown {
				mb.Shutdown()
			}
			
			if got := mb.IsShutdown(); got != tt.want {
				t.Errorf("Microbatch.IsShutdown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMicrobatch_Shutdown(t *testing.T) {
	type fields struct {
		RWMutex    sync.RWMutex
		size       int
		maxAge     time.Duration
		proc       BatchProcessor
		queue      chan Batch
		submit     chan Job
		Results    map[int]JobResult
		shutdown   chan struct{}
		numWorkers int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Microbatch{
				RWMutex:    tt.fields.RWMutex,
				size:       tt.fields.size,
				maxAge:     tt.fields.maxAge,
				proc:       tt.fields.proc,
				queue:      tt.fields.queue,
				submit:     tt.fields.submit,
				Results:    tt.fields.Results,
				shutdown:   tt.fields.shutdown,
				numWorkers: tt.fields.numWorkers,
			}
			b.Shutdown()
		})
	}
}

func TestMicrobatch_StoreJobResults(t *testing.T) {
	type fields struct {
		RWMutex    sync.RWMutex
		size       int
		maxAge     time.Duration
		proc       BatchProcessor
		queue      chan Batch
		submit     chan Job
		Results    map[int]JobResult
		shutdown   chan struct{}
		numWorkers int
	}
	type args struct {
		jr []JobResult
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Microbatch{
				RWMutex:    tt.fields.RWMutex,
				size:       tt.fields.size,
				maxAge:     tt.fields.maxAge,
				proc:       tt.fields.proc,
				queue:      tt.fields.queue,
				submit:     tt.fields.submit,
				Results:    tt.fields.Results,
				shutdown:   tt.fields.shutdown,
				numWorkers: tt.fields.numWorkers,
			}
			b.StoreJobResults(tt.args.jr)
		})
	}
}

func TestMicrobatch_PrintAllJobResults(t *testing.T) {
	type fields struct {
		RWMutex    sync.RWMutex
		size       int
		maxAge     time.Duration
		proc       BatchProcessor
		queue      chan Batch
		submit     chan Job
		Results    map[int]JobResult
		shutdown   chan struct{}
		numWorkers int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Microbatch{
				RWMutex:    tt.fields.RWMutex,
				size:       tt.fields.size,
				maxAge:     tt.fields.maxAge,
				proc:       tt.fields.proc,
				queue:      tt.fields.queue,
				submit:     tt.fields.submit,
				Results:    tt.fields.Results,
				shutdown:   tt.fields.shutdown,
				numWorkers: tt.fields.numWorkers,
			}
			b.PrintAllJobResults()
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		size       int
		maxAge     time.Duration
		numWorkers int
		proc       BatchProcessor
	}
	tests := []struct {
		name string
		args args
		want Microbatch
	}{
		// TODO: Add test cases.
	}
 	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.size, tt.args.maxAge, tt.args.numWorkers, tt.args.proc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
