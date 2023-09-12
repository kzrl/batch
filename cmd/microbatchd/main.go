package main

import(
	"github.com/kzrl/microbatch"
	"time"
	"log/slog"
)


type MyWork struct{}
func (m *MyWork) Do() (error, string) {
	msg := "This is the result of Do() defined in the program calling the library"
	slog.Info("MyWork.Do()", "output", msg)
	return nil, msg
}

func main() {
	// NoopBatchProcess is an example implementation of BatchProcessor interface.
	var nbp microbatch.NoopBatchProcessor

	// New Microbatch with batches of 10, maxAge 10 seconds, using the NoopBatchProcessor.
	mb := microbatch.New(10, time.Second * 10, 2, &nbp) 

	// Start the SubmitWorkers.
	mb.Start()
	// Create 42 new jobs, each with a MyWork.
	for i:=0; i < 42; i++ { 
		j := microbatch.NewJob(i)
		var m MyWork
		j.Work = &m
		_ = mb.Submit(j) //ignoring the returned JobResult because it's not interesting yet.
	}
	// Signal to the workers that it's time to finish up and exit.
	mb.Shutdown()

	// We've shutdown already, so this doesn't do anything, it's just not processed.
	j := microbatch.NewJob(5000)
	mb.Submit(j)

	mb.PrintAllJobResults()
}
