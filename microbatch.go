package microbatch

import(
	"time"
	"sync"
	"fmt"
)


// Microbatch holds the configuration for desired batch size and frequency.
type Microbatch struct {
	sync.RWMutex
	size int
	maxAge time.Duration
	jobs []Job
	proc BatchProcessor
	queue chan(Batch)
	submitCh chan(Job)
	allDone chan struct{}
}

func (b *Microbatch) Submit(j Job) JobResult {
	fmt.Printf("Submit job: %s\n", &j)
	b.Lock()	
	defer b.Unlock()	
	b.jobs = append(b.jobs, j)
	fmt.Printf("Added to slice: %s\n", &j)
	var jr JobResult
	jr.Output = j.String()
	return jr
}

func (b *Microbatch) CountJobs() int {
	b.RLock()
	defer b.RUnlock()
	return len(b.jobs)
}

func (b *Microbatch) buildJobs()  {
	
	for jobs := range submitCh {

	}
}

func (b *Microbatch) Shutdown() {
	
	//<-b.allDone
}

func (b *Microbatch) Process() {
	go b.proc.Do(b.queue)
	
	for {
		b.RLock()
		defer b.RUnlock()
		if len(b.jobs) == 0 {
			fmt.Println("No jobs to process")
			time.Sleep(5 * time.Second) //TODO.
		}		
		// check if we have enough jobs for a batch
		if len(b.jobs) >= b.size {
			fmt.Println("Got enough jobs to make a batch")
			thisBatch := b.jobs[0:b.size]
			fmt.Printf("len(thisMicrobatch) = %d\n", len(thisBatch))
			time.Sleep(5 * time.Second) // included so we can see this working more easily.
			remaining := b.jobs[b.size:]
			b.jobs = remaining

			// make a Batch and send it to the queue.
			batch := Batch{Jobs: thisBatch}
			b.queue <- batch
			
			
			
		} else {
			// check if the jobs are old enough to be sent in a smaller batch.		
			notOldEnough := make([]Job, 0)
			for _, job := range b.jobs {
				if time.Since(job.created) >= b.maxAge {
					// if this job is old enough, send it anyway.
					fmt.Printf("job is old enough, send it\n")
					b.queue <- Batch{Jobs: []Job{job}}
				} else {
					//fmt.Printf("%s not old enough\n", &job)
					// otherwise append it to notOldEnough slice.
					notOldEnough = append(notOldEnough, job)
				}
			}
			// jobs left to process on next iteration are the remaining ones.
			b.jobs = notOldEnough
		}
	}
}


// BatchProcessor processes jobs and returns a JobResult for each.
type BatchProcessor interface {
	Do(<- chan Batch) []JobResult
}

type NoopBatchProcessor struct{}

func (nbp *NoopBatchProcessor) Do(batches <- chan(Batch)) []JobResult {
	for batch := range batches {
		for _, job := range batch.Jobs {
			fmt.Printf("received: %s\n", &job)
		fmt.Println(&job)
		}
		
	}
	fmt.Println("DONE")
	return nil;
}



// New returns a configured Batch ready for use.
func New(size int, maxAge time.Duration, proc BatchProcessor) Microbatch {
	var b Microbatch
	b.size = size
	b.maxAge = maxAge
	b.proc = proc
	b.allDone = make(chan struct{})
	b.queue = make(chan Batch, 0) // send Jobs to the processor in an unbuffered channel.

	return b
}

