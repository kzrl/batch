package microbatch

import (
	"fmt"
	"time"
)

// JobResult contains the output of a Job or an error.
type JobResult struct {
	Output string
	Job
	Error error
}

func (jr *JobResult) String() string {
	return fmt.Sprintf("Job.IDD=%d created=%s completed=%t)", jr.ID, jr.created, jr.Job.IsDone())
}

// Job is executed by a BatchProcessor.
type Job struct {
	ID        int
	created   time.Time
	completed time.Time //TODO - these should be private to avoid other things mutating them at runtime.
	Work
}

// Work is interface with a single method Do().
// Calling programs can implement this interface
// to do the desired operations..
type Work interface {
	Do() (error, string)
}

// IsDone checks if the Job is marked as completed.
func (j *Job) IsDone() bool {
	return !j.completed.IsZero()
}

func (j *Job) String() string {
	return fmt.Sprintf("ID=%d created=%s completed=%t)", j.ID, j.created, j.IsDone())
}

func (j *Job) Created() time.Time {
	return j.created
}

// Done marks the Job as completed.
func (j *Job) Done() error {
	j.completed = time.Now()
	return nil
}

func NewJob(id int) Job {
	return Job{
		ID:      id,
		created: time.Now(),
	}
}

type Batch struct {
	Jobs []Job
}
