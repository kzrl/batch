package main

import(
	"github.com/kzrl/microbatch"
	"time"
	"sync"
)



func main() {

	var nbp microbatch.NoopBatchProcessor
	var wg sync.WaitGroup
	mb := microbatch.New(50, time.Second * 10, &nbp)
	
	for i:=0; i< 20; i++ {
		j := microbatch.NewJob(i)
		wg.Add(1)

		// Need the function literal to avoid main() finishing
		// before Submit() has unlocked, it was causing deadlocks.
		go func() {
			defer wg.Done()
			mb.Submit(j)
		}()
		
	}
	mb.Process()
	wg.Wait()

	//batch.Shutdown()
	

}
