package microbatch

import (
	"reflect"
	"testing"
	"time"
)

func TestJob_IsDone(t *testing.T) {

	tests := []struct {
		name string
		done bool
		want bool
	}{
		{name: "Job is marked as done", done: true, want: true},
		{name: "Job is not marked as done", done: false, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := NewJob(42)
			if tt.done {
				j.Done()
			}
			if got := j.IsDone(); got != tt.want {
				t.Errorf("Job.IsDone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_String(t *testing.T) {
	type fields struct {
		ID        int
		created   time.Time
		completed time.Time
		Work      Work
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				ID:        tt.fields.ID,
				created:   tt.fields.created,
				completed: tt.fields.completed,
				Work:      tt.fields.Work,
			}
			if got := j.String(); got != tt.want {
				t.Errorf("Job.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_Done(t *testing.T) {
	type fields struct {
		ID        int
		created   time.Time
		completed time.Time
		Work      Work
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				ID:        tt.fields.ID,
				created:   tt.fields.created,
				completed: tt.fields.completed,
				Work:      tt.fields.Work,
			}
			if err := j.Done(); (err != nil) != tt.wantErr {
				t.Errorf("Job.Done() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewJob(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name string
		args args
		want Job
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJob(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJob() = %v, want %v", got, tt.want)
			}
		})
	}
}
