package microbatch

import(
	"time"
	"fmt"
)


type JobResult struct {
	Output string
}

type Work interface {
	Do() JobResult
}

type Job struct {
	id int
	created time.Time
	completed time.Time
}

func (j *Job) IsDone() bool {
	return !j.completed.IsZero()
}

func (j *Job) String() string {
	return fmt.Sprintf("id=%d created=%s completed=%t)", j.id, j.created, j.IsDone())
}

func NewJob(id int) Job {
	return Job{
		id: id,
		created: time.Now(),
	}
}

type Batch struct {
	Jobs []Job
}
