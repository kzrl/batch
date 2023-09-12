package microbatch

import (
	"reflect"
	"testing"
)

func TestNoopBatchProcessor_Do(t *testing.T) {
	type args struct {
		batches <-chan (Batch)
	}
	tests := []struct {
		name string
		nbp  *NoopBatchProcessor
		args args
		want []JobResult
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nbp := &NoopBatchProcessor{}
			if got := nbp.Do(tt.args.batches); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NoopBatchProcessor.Do() = %v, want %v", got, tt.want)
			}
		})
	}
}
