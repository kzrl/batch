// Package microbatch provides a way to create Batches of Jobs.
package microbatch

import (
	"log/slog"
	"sync"
	"time"
)

// buffSize is the size of the buffered submit and queue channels.
// Can tune this based on application requirements.
const buffSize = 100

// Microbatch holds the configuration for desired batch size and frequency.
type Microbatch struct {
	sync.RWMutex
	size       int
	maxAge     time.Duration
	proc       BatchProcessor
	queue      chan Batch
	submit     chan Job
	results    sync.Map
	shutdown   chan struct{}
	once       sync.Once // used to avoid a closed channel panic if we call .Shutdown() twice.
	numWorkers int
}

var wg sync.WaitGroup

// Submit adds a job to the submit channel and returns a JobResult.
func (b *Microbatch) Submit(j Job) JobResult {
	slog.Info("Submitted", "job.ID", j.ID, "created", j.Created())
	b.submit <- j
	var jr JobResult // return a JobResult which includes the Job and some details in the Output.
	jr.Job = j
	jr.Output = j.String()
	return jr
}

// Start starts the configured number of workers. Returns the number of goroutines started.
func (b *Microbatch) Start() int {
	for i := 0; i < b.numWorkers; i++ {
		wg.Add(1)
		go b.submitWorker(i)
	}
	return b.numWorkers
}

// submitWorker reads off submit channel and builds the batches of jobs.
func (b *Microbatch) submitWorker(id int) {
	defer wg.Done()

	slog.Info("Start SubmitWorker", "id", id)
	receivedCount := 0
	processedCount := 0
	jobs := make([]Job, 0) // to avoid needing mutexes, each worker gets their own slice of jobs.
	for {
		// shutdown has been signalled and we're out of jobs to process, return.
		if b.IsShutdown() && len(jobs) == 0 {
			slog.Info("SubmitWorker finished all jobs", "id", id, "received", receivedCount, "processed", processedCount)
			return
		}
		if len(jobs) == b.size {
			slog.Info("Made a batch", "size", b.size)
			batch := Batch{Jobs: jobs}
			processedCount += len(jobs)
			b.queue <- batch
			jr := b.proc.Do(b.queue)
			b.StoreJobResults(jr)
			jobs = make([]Job, 0) // reset the jobs in flight.
		}

		for _, job := range jobs {
			// if we don't have enough jobs to make a batch, but one of the in-flight jobs is old enough, we can send them all.
			if time.Since(job.Created()) >= b.maxAge {
				batch := Batch{Jobs: jobs}
				b.queue <- batch
				jr := b.proc.Do(b.queue)
				b.StoreJobResults(jr)
				processedCount += len(jobs)
				jobs = make([]Job, 0) // reset the jobs in flight.
			}
		}
		select {
		case job := <-b.submit:
			slog.Info("SubmitWorker", "workerID", id, "job.ID", job.ID, "job.Created", job.Created)
			jobs = append(jobs, job)
			receivedCount++
		default:
			// default case so that select doesn't block and we keep looping.
		}
	}
}

// IsShutDown helper method to listen for shutdown signal.
func (b *Microbatch) IsShutdown() bool {
	select {
	case <-b.shutdown:
		return true
	default:
		return false
	}
}

// Shutdown closes the b.shutdown channel to signal to SubmitWorker goroutines to finish work and exit.
func (b *Microbatch) Shutdown() {
	slog.Info("Microbatch shutting down")
	b.once.Do(func() {
		close(b.shutdown) // close the shutdown channel to signal to SubmitWorkers to finish up their work.
	})
	wg.Wait()
	slog.Info("Microbatch shut down")
}

// StoreJobResults populates the b.Results map.
// Uses a sync.Map to make it safe to be called from multiple goroutines.
func (b *Microbatch) StoreJobResults(jr []JobResult) {
	for _, result := range jr {
		slog.Info("StoreJobResults", "id", result.Job.ID)
		b.results.Store(result.Job.ID, result)
	}

}

// PrintAllJobResults loops over the b.Results sync.Map.
func (b *Microbatch) PrintAllJobResults() {
	b.results.Range(func(id any, result any) bool {
		jr, ok := result.(JobResult)
		if ok {
			slog.Info("JobResult[id]", "id", jr.ID, "output", jr.Output, "error", jr.Error, "IsDone()", jr.Job.IsDone())
		}

		return true
	})
}

// GetJobResult retrieves a JobResult from the b.results sync.Map.
func (b *Microbatch) GetJobResult(id int) (JobResult, bool) {
	result, foundOk := b.results.Load(id)
	if !foundOk {
		return JobResult{}, false
	}
	jr, ok := result.(JobResult)
	if !ok {
		slog.Error("GetJobResult(id) failed type assertion to JobResult", "id", id)
	}
	return jr, ok
}

// New returns a configured Microbatch ready for use.
func New(size int, maxAge time.Duration, numWorkers int, proc BatchProcessor) Microbatch {
	var b Microbatch
	b.size = size
	b.maxAge = maxAge
	b.proc = proc
	b.shutdown = make(chan struct{})
	b.queue = make(chan Batch, buffSize)
	b.submit = make(chan Job, buffSize)
	b.numWorkers = numWorkers
	return b
}
