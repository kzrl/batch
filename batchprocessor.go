package microbatch

import(
	"log/slog"
)

// BatchProcessor processes jobs and returns a JobResult for each.
type BatchProcessor interface {
	Do(<- chan Batch) []JobResult
}

// NoopBatchProcessor is an example implementation of BatchProcessor interface.
type NoopBatchProcessor struct{}

// Do receives Batch types from the channel argument.
// For each Job in a batch, it executes the job, marks it as done.
func (nbp *NoopBatchProcessor) Do(batches <- chan(Batch)) []JobResult {
	var jr []JobResult
	select {
	case batch := <- batches:
		for _, job := range batch.Jobs {
			// nil guard in case job.Work isn't set.
			if job.Work == nil {
				job.Done()
				jr = append(jr, JobResult{Job: job})
				continue
			}
			err, result := job.Do()
			// not checking err != nil, just adding it to the JobResult.
			job.Done()
			slog.Info("NoopBatchProcessor", "job.ID", job.ID, "job.IsDone", job.IsDone())
			jr = append(jr, JobResult{Output: result, Job: job, Error: err})
		}
	}
	return jr;
}

