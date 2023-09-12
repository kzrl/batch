// Package microbatch provides a way to create Batches of Jobs.
package microbatch

import(
	"time"
	"sync"
	"log/slog"
)


// buffSize is the size of the buffered submit and queue channels.
// Can tune this based on application requirements.
const buffSize = 100

// Microbatch holds the configuration for desired batch size and frequency.
type Microbatch struct {
	sync.RWMutex
	size int
	maxAge time.Duration
	proc BatchProcessor
	queue chan Batch
	submit chan Job
	Results map[int]JobResult
	shutdown chan struct{}
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
	for i:= 0; i< b.numWorkers; i++ {
		wg.Add(1)
		go b.SubmitWorker(i)
	}
	return b.numWorkers
}

// SubmitWorker reads off submit channel and builds the batches of jobs.
func (b *Microbatch) SubmitWorker(id int) {
	//TODO - bug, if we call SubmitWorker directly without using Start(), this will panic with negative WaitGroup counter
	defer wg.Done() 

	slog.Info("Start SubmitWorker", "id", id)
	receivedCount := 0
	processedCount := 0
	jobs := make([]Job, 0) // to avoid needing mutexes, each worker gets their own slice of jobs
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
	close(b.shutdown) // close the shutdown channel to signal to SubmitWorkers to finish up their work.
	wg.Wait()
}

// StoreJobResults populates the b.Results map.
// Uses an RWMutex.Lock() to make it safe to be called from multiple goroutines.
func (b *Microbatch) StoreJobResults(jr []JobResult) {
	b.Lock()
	for _, result := range jr {
		b.Results[result.Job.ID] = result
	}
	b.Unlock()
}

// PrintAllJobResults loops over the b.Results map.
// Uses an RWMutex.RLock to make it safe to be called from multiple goroutines.
func (b *Microbatch) PrintAllJobResults() {
	b.RLock()
	for id, result := range b.Results {
		slog.Info("JobResult[id]", "id", id, "output", result.Output, "error", result.Error, "IsDone()", result.Job.IsDone())
	}
	b.RUnlock()
}




// New returns a configured Microbatch ready for use.
func New(size int, maxAge time.Duration, numWorkers int, proc BatchProcessor) Microbatch {
	var b Microbatch
	b.size = size
	b.maxAge = maxAge
	b.proc = proc
	b.shutdown = make(chan struct{})
	b.queue = make(chan Batch, buffSize) // buffer size chosen at random.
	b.submit = make(chan Job, buffSize) // buffer size chosen at random.
	b.numWorkers = numWorkers
	b.Results = make(map[int]JobResult)
	return b
}

