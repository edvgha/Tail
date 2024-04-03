package space

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"os"
	"tail.server/app/optimizer/code/misc"
	"testing"
)

func TestSpace_sampleBuckets(t *testing.T) {
	type args struct {
		price float64
	}
	space, err := NewSpace("ctx", 0.2, 10.0, misc.Config{LevelSize: 3, BucketSize: 30, BufferSize: 10, Discount: 0.25, DesiredExplorationSpeed: 2}, zerolog.New(os.Stdout))
	assert.Nil(t, err)

	tests := []struct {
		name  string
		space *Space
		args  args
		want  int
	}{
		{"sampleBuckets_1", space, args{2.2}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, len(space.sampleBuckets(tt.args.price)), "sampleBuckets(%v)", tt.args.price)
		})
	}
}
