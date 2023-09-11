package microbatch

import(
	"time"
	"fmt"
)


type JobResult struct {
	Output string
}

type Job struct {
	ID int
	Created time.Time
	Completed time.Time
}

func (j *Job) IsDone() bool {
	return !j.Completed.IsZero()
}

func (j *Job) String() string {
	return fmt.Sprintf("id=%d created=%s completed=%t)", j.ID, j.Created, j.IsDone())
}

func NewJob(id int) Job {
	return Job{
		ID: id,
		Created: time.Now(),
	}
}

type Batch struct {
	Jobs []Job
}
