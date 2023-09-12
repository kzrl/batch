# microbatch

A Go library to implement [microbatch processing](https://hazelcast.com/glossary/micro-batch-processing/)

Note: This is a technical test, not intended for production use as-is.


### Task requirements

 Micro-batching is a technique used in processing pipelines where individual tasks are grouped
 together into small batches. This can improve throughput by reducing the number of requests made
 to a downstream system. Your task is to implement a micro-batching library, with the following
 requirements:
- it should allow the caller to submit a single Job, and it should return a JobResult [x]
- it should process accepted Jobs in batches using a BatchProcessor [x]
- Don't implement BatchProcessor. This should be a dependency of your library. Note: example implementation included.
- it should provide a way to configure the batching behaviour i.e. size and frequency [x]
- it should expose a shutdown method which returns after all previously accepted Jobs are processed [x]


## Library usage

Example in cmd/microbatchd/main.go

```golang
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
```

Example output:

``` 
~/src/microbatch$ go run -race cmd/microbatchd/main.go   
2023/09/12 09:39:30 INFO Submitted job.ID=0 created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO Submitted job.ID=1 created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO Start SubmitWorker id=0
2023/09/12 09:39:30 INFO Start SubmitWorker id=1
2023/09/12 09:39:30 INFO Submitted job.ID=2 created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=1 job.ID=1 job.Created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=0 job.ID=0 job.Created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO Submitted job.ID=3 created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=0 job.ID=3 job.Created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=1 job.ID=2 job.Created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO Submitted job.ID=4 created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=1 job.ID=4 job.Created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO Submitted job.ID=5 created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=0 job.ID=5 job.Created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO Submitted job.ID=6 created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=0 job.ID=6 job.Created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO Submitted job.ID=7 created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=1 job.ID=7 job.Created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO Submitted job.ID=8 created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=0 job.ID=8 job.Created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO Submitted job.ID=9 created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO Submitted job.ID=10 created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=1 job.ID=9 job.Created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=0 job.ID=10 job.Created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO Submitted job.ID=11 created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=1 job.ID=11 job.Created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO Submitted job.ID=12 created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=0 job.ID=12 job.Created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO Submitted job.ID=13 created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=1 job.ID=13 job.Created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO Submitted job.ID=14 created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO Submitted job.ID=15 created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=0 job.ID=14 job.Created=2023-09-12T09:39:30.527+10:00
2023/09/12 09:39:30 INFO Submitted job.ID=16 created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=1 job.ID=15 job.Created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO Submitted job.ID=17 created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=0 job.ID=16 job.Created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO Submitted job.ID=18 created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=1 job.ID=17 job.Created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO Submitted job.ID=19 created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=0 job.ID=18 job.Created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO Submitted job.ID=20 created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO Made a batch size=10
2023/09/12 09:39:30 INFO SubmitWorker workerID=1 job.ID=19 job.Created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO Made a batch size=10
2023/09/12 09:39:30 INFO Submitted job.ID=21 created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO Submitted job.ID=22 created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=0 job.IsDone=true
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=1 job.IsDone=true
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO Submitted job.ID=23 created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=3 job.IsDone=true
2023/09/12 09:39:30 INFO Submitted job.ID=24 created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=2 job.IsDone=true
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO Submitted job.ID=25 created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=4 job.IsDone=true
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=5 job.IsDone=true
2023/09/12 09:39:30 INFO Submitted job.ID=26 created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO Submitted job.ID=27 created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=7 job.IsDone=true
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=6 job.IsDone=true
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO Submitted job.ID=28 created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=9 job.IsDone=true
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO Submitted job.ID=29 created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=8 job.IsDone=true
2023/09/12 09:39:30 INFO Submitted job.ID=30 created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=11 job.IsDone=true
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO Submitted job.ID=31 created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=10 job.IsDone=true
2023/09/12 09:39:30 INFO Submitted job.ID=32 created=2023-09-12T09:39:30.529+10:00
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=13 job.IsDone=true
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO Submitted job.ID=33 created=2023-09-12T09:39:30.529+10:00
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=12 job.IsDone=true
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=15 job.IsDone=true
2023/09/12 09:39:30 INFO Submitted job.ID=34 created=2023-09-12T09:39:30.529+10:00
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO Submitted job.ID=35 created=2023-09-12T09:39:30.529+10:00
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=14 job.IsDone=true
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=17 job.IsDone=true
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO Submitted job.ID=36 created=2023-09-12T09:39:30.529+10:00
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=16 job.IsDone=true
2023/09/12 09:39:30 INFO Submitted job.ID=37 created=2023-09-12T09:39:30.529+10:00
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=19 job.IsDone=true
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO SubmitWorker workerID=1 job.ID=20 job.Created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO Submitted job.ID=38 created=2023-09-12T09:39:30.529+10:00
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=18 job.IsDone=true
2023/09/12 09:39:30 INFO SubmitWorker workerID=1 job.ID=21 job.Created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO Submitted job.ID=39 created=2023-09-12T09:39:30.529+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=0 job.ID=22 job.Created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=1 job.ID=23 job.Created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO Submitted job.ID=40 created=2023-09-12T09:39:30.529+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=0 job.ID=24 job.Created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=1 job.ID=25 job.Created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO Submitted job.ID=41 created=2023-09-12T09:39:30.529+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=0 job.ID=26 job.Created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO Microbatch shutting down
2023/09/12 09:39:30 INFO SubmitWorker workerID=1 job.ID=27 job.Created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=0 job.ID=28 job.Created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=1 job.ID=29 job.Created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=0 job.ID=30 job.Created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=1 job.ID=31 job.Created=2023-09-12T09:39:30.528+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=0 job.ID=32 job.Created=2023-09-12T09:39:30.529+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=1 job.ID=33 job.Created=2023-09-12T09:39:30.529+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=0 job.ID=34 job.Created=2023-09-12T09:39:30.529+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=1 job.ID=35 job.Created=2023-09-12T09:39:30.529+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=0 job.ID=36 job.Created=2023-09-12T09:39:30.529+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=1 job.ID=37 job.Created=2023-09-12T09:39:30.529+10:00
2023/09/12 09:39:30 INFO SubmitWorker workerID=0 job.ID=38 job.Created=2023-09-12T09:39:30.529+10:00
2023/09/12 09:39:30 INFO Made a batch size=10
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO SubmitWorker workerID=0 job.ID=39 job.Created=2023-09-12T09:39:30.529+10:00
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=20 job.IsDone=true
2023/09/12 09:39:30 INFO Made a batch size=10
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=21 job.IsDone=true
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=22 job.IsDone=true
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=23 job.IsDone=true
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=24 job.IsDone=true
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=25 job.IsDone=true
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=26 job.IsDone=true
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=27 job.IsDone=true
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=28 job.IsDone=true
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=29 job.IsDone=true
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=31 job.IsDone=true
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=30 job.IsDone=true
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=33 job.IsDone=true
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=32 job.IsDone=true
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=35 job.IsDone=true
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=34 job.IsDone=true
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=37 job.IsDone=true
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=36 job.IsDone=true
2023/09/12 09:39:30 INFO SubmitWorker workerID=1 job.ID=40 job.Created=2023-09-12T09:39:30.529+10:00
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO SubmitWorker workerID=1 job.ID=41 job.Created=2023-09-12T09:39:30.529+10:00
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=38 job.IsDone=true
2023/09/12 09:39:30 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:30 INFO NoopBatchProcessor job.ID=39 job.IsDone=true
2023/09/12 09:39:30 INFO SubmitWorker finished all jobs id=0 received=20 processed=20
2023/09/12 09:39:40 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:40 INFO NoopBatchProcessor job.ID=40 job.IsDone=true
2023/09/12 09:39:40 INFO MyWork.Do() output="This is the result of Do() defined in the program calling the library"
2023/09/12 09:39:40 INFO NoopBatchProcessor job.ID=41 job.IsDone=true
2023/09/12 09:39:40 INFO SubmitWorker finished all jobs id=1 received=22 processed=22
2023/09/12 09:39:40 INFO Submitted job.ID=5000 created=2023-09-12T09:39:40.529+10:00
2023/09/12 09:39:40 INFO JobResult[id] id=21 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=34 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=35 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=22 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=7 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=19 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=0 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=8 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=29 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=33 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=32 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=39 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=9 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=13 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=12 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=14 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=25 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=40 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=1 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=2 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=15 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=6 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=10 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=24 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=28 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=4 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=11 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=17 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=18 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=23 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=5 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=20 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=31 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=26 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=41 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=36 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=38 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=3 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=16 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=27 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=37 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true
2023/09/12 09:39:40 INFO JobResult[id] id=30 output="This is the result of Do() defined in the program calling the library" error=<nil> IsDone()=true

```
