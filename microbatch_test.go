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
		name string
		job  Job
		want JobResult
	}{
		{name: "Submit a good job", job: testJob, want: JobResult{Job: testJob, Output: testJob.String()}},
		{name: "Submit an empty job", job: emptyJob, want: JobResult{Job: emptyJob, Output: emptyJob.String()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var nbp NoopBatchProcessor
			mb := New(batchSize, maxAge, numWorkers, &nbp)

			if got := mb.Submit(tt.job); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Microbatch.Submit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMicrobatch_Start(t *testing.T) {
	tests := []struct {
		name       string
		want       int
		numWorkers int
	}{
		{name: "Test Start() with 2 workers", numWorkers: 2, want: 2},
		{name: "Test Start() with -1 workers", numWorkers: -1, want: 2},
		{name: "Test Start() with 20 workers", numWorkers: 20, want: 20},
		{name: "Test Start() with 0 workers", numWorkers: 0, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var nbp NoopBatchProcessor
			mb := New(batchSize, maxAge, tt.numWorkers, &nbp)
			got := mb.Start()
			if got != tt.numWorkers {
				t.Errorf("Microbatch.Start() = %v, want %v", got, tt.want)
			}
			mb.Shutdown()
		})
	}
}

func TestMicrobatch_submitWorker(t *testing.T) {
	tests := []struct {
		name       string
		want       int
		numWorkers int
		numJobs    int
	}{
		{name: "Check that submitWorker processes submitted jobs OK", numWorkers: 2, numJobs: 42},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var nbp NoopBatchProcessor
			mb := New(batchSize, maxAge, tt.numWorkers, &nbp)
			got := mb.Start()
			if got != tt.numWorkers {
				t.Errorf("Microbatch.Start() = %v, want %v", got, tt.want)
			}

			// Submit some jobs.
			for i := 0; i < tt.numJobs; i++ {
				j := NewJob(i)
				mb.Submit(j)
			}
			mb.Shutdown()
			for i := 0; i < tt.numJobs; i++ {
				_, found := mb.GetJobResult(i)
				if !found {
					t.Errorf("submitWorker didn't store job correctly. GetJobResult(%d) = _, %t", i, found)
				}
			}
		})
	}
}

func TestMicrobatch_IsShutdown(t *testing.T) {
	var numWorkers = 2
	tests := []struct {
		name     string
		start    bool
		shutdown bool
		want     bool
	}{
		{name: "Start it up and shut it down again", start: true, shutdown: true, want: true},
		{name: "Start it up and leave it running", start: true, shutdown: false, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var nbp NoopBatchProcessor
			mb := New(batchSize, maxAge, numWorkers, &nbp)
			if tt.start {
				mb.Start()
			}
			if tt.shutdown {
				mb.Shutdown()
			}

			if got := mb.IsShutdown(); got != tt.want {
				t.Errorf("Microbatch.IsShutdown() = %v, want %v", got, tt.want)
			}
			mb.Shutdown()
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
		results    map[int]JobResult
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
				shutdown:   tt.fields.shutdown,
				numWorkers: tt.fields.numWorkers,
			}
			b.Shutdown()
		})
	}
}

func TestMicrobatch_StoreJobresults(t *testing.T) {

	type args struct {
		jr []JobResult
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestMicrobatch_PrintAllJobresults(t *testing.T) {
	type fields struct {
		RWMutex    sync.RWMutex
		size       int
		maxAge     time.Duration
		proc       BatchProcessor
		queue      chan Batch
		submit     chan Job
		results    map[int]JobResult
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

func TestMicrobatch_StoreJobResults(t *testing.T) {
	type fields struct {
		RWMutex    sync.RWMutex
		size       int
		maxAge     time.Duration
		proc       BatchProcessor
		queue      chan Batch
		submit     chan Job
		results    map[int]JobResult
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
		results    map[int]JobResult
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

		})
	}
}

func TestMicrobatch_GetJobResult(t *testing.T) {

	numWorkers := 2

	tests := []struct {
		name string
		id   int
		got  JobResult
		ok   bool
	}{
		{name: "Got JobResult from an ID", id: 2, ok: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var nbp NoopBatchProcessor
			mb := New(batchSize, maxAge, numWorkers, &nbp)
			mb.Start()

			// Submit some jobs.
			for i := 0; i < 42; i++ {
				j := NewJob(i)
				mb.Submit(j)
			}
			mb.Shutdown()
			for i := 0; i < 42; i++ {
				_, found := mb.GetJobResult(i)
				if !found {
					t.Errorf("GetJobResult(%d) = _, %t", i, found)
				}
			}
		})
	}
}
