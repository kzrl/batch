package microbatch

import (
	"reflect"
	"testing"
	"time"
)

func TestBatch_Submit(t *testing.T) {
	type fields struct {
		size int
		freq time.Duration
	}
	type args struct {
		j Job
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   JobResult
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Batch{
				size: tt.fields.size,
				freq: tt.fields.freq,
			}
			if got := b.Submit(tt.args.j); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Batch.Submit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		size int
		freq time.Duration
	}
	tests := []struct {
		name string
		args args
		want Batch
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.size, tt.args.freq); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
