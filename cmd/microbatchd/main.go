package main

import(
	"github.com/kzrl/microbatch"
	"time"
)



func main() {
	var nbp microbatch.NoopBatchProcessor

	mb := microbatch.New(10, time.Second * 10, 2, &nbp)
	
	for i:=0; i< 100; i++ {
		j := microbatch.NewJob(i)
		mb.Submit(j)	
	}
	mb.Shutdown()
}
