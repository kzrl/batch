package microbatch

import(
	"time"
	"sync"
	"fmt"
)



// Microbatch holds the configuration for desired batch size and frequency.
type Microbatch struct {
	size int
	maxAge time.Duration
	proc BatchProcessor
	queue chan Batch
	submit chan Job
	shutdown chan struct{}
	numWorkers int
}

var wg sync.WaitGroup

func (b *Microbatch) Submit(j Job) JobResult {
	fmt.Printf("Submitted job: %s\n", &j)
	b.submit <- j

	//TODO
	var jr JobResult
	jr.Output = j.String()
	return jr
}

// Start creates the configured number of workers.
func (b *Microbatch) Start() {
	for i:= 0; i< b.numWorkers; i++ {
		go b.SubmitWorker(i)
	}
}

// SubmitWorker reads off submit channel and builds the batches of jobs.
func (b *Microbatch) SubmitWorker(id int) {
	fmt.Printf("Start worker %d\n", id)
	jobs := make([]Job, 0) // to avoid needing mutexes, each worker gets their own slice of jobs
 	for {
		// shutdown signalled and we're out of jobs to process, return.
		if b.IsShutdown() && len(jobs) == 0 {
			fmt.Println("Worker all done")
			return
		}
		if len(jobs) == b.size {
			fmt.Println("Got a whole batch")
			batch := Batch{Jobs: jobs}
			b.queue <- batch
			b.proc.Do(b.queue)
			jobs = make([]Job, 0) // reset the jobs in flight.
		}

		for _, job := range jobs {
			// if we don't have enough jobs to make a batch, but one of the in-flight jobs is old enough, we can send them all.
			if time.Since(job.Created) >= b.maxAge {
				batch := Batch{Jobs: jobs}
				b.queue <- batch
				b.proc.Do(b.queue)
				jobs = make([]Job, 0) // reset the jobs in flight.
			}
		}
		select {
		case job := <-b.submit:
			fmt.Printf("Worker %d received job: %s\n", id, &job)
			jobs = append(jobs, job)
		default:
			// default case so that select doesn't block and we keep looping.
		}
	}
}

func (b *Microbatch) IsShutdown() bool {
	select {
	case <-b.shutdown:
		return true
	default:
		return false
	}
}

func (b *Microbatch) Shutdown() {
	close(b.shutdown) // close the shutdown channel to signal to SubmitWorkers to finish up.
}

// BatchProcessor processes jobs and returns a JobResult for each.
type BatchProcessor interface {
	Do(<- chan Batch) []JobResult
}

type NoopBatchProcessor struct{}

func (nbp *NoopBatchProcessor) Do(batches <- chan(Batch)) []JobResult {

	select {
	case batch := <- batches:
		for _, job := range batch.Jobs {
			fmt.Printf("NoopBatchProcessor received: %s\n", &job)
		}
	}
	fmt.Println("NoopBatchProcessor Done")
	return nil;
}



// New returns a configured Batch ready for use.
func New(size int, maxAge time.Duration, numWorkers int, proc BatchProcessor) Microbatch {
	var b Microbatch
	b.size = size
	b.maxAge = maxAge
	b.proc = proc
	b.shutdown = make(chan struct{})
	b.queue = make(chan Batch, 100) //buffer size chosen at random.
	b.submit = make(chan Job, 100) //buffer size chosen at random.
	b.numWorkers = numWorkers
	b.Start() //start the submit workers.
	return b
}

