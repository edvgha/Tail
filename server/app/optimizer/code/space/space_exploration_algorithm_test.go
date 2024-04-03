package space

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestInitUniformFlat(t *testing.T) {
	type args struct {
		context          string
		minPrice         float64
		maxPrice         float64
		nBins            int
		desiredFrequency float64
	}
	tests := []struct {
		name string
		args args
		want UniformFlat
	}{
		{"InitUniformFlat_0",
			args{"",
				0.1,
				0.5,
				4,
				5},
			UniformFlat{
				bins: []float64{0.1, 0.2, 0.30000000000000004, 0.4, 0.5},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iuf := InitUniformFlat(tt.args.context, tt.args.minPrice, tt.args.maxPrice, tt.args.nBins, tt.args.desiredFrequency, zerolog.New(os.Stdout))
			assert.True(t, reflect.DeepEqual(iuf.bins, tt.want.bins))
		})
	}
}

func TestUniformFlat_findLeftmost(t *testing.T) {
	type fields struct {
		bins         []float64
		context      string
		lastExplored []time.Time
		desiredSpeed float64
	}
	type args struct {
		price float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  bool
	}{
		{"findLeftMost_0", fields{
			bins:         []float64{0.1, 0.2, 0.3, 0.4, 0.5},
			context:      "",
			lastExplored: make([]time.Time, 0),
			desiredSpeed: 0.25},
			args{0.22}, 1, true},
		{"findLeftMost_1", fields{
			bins:         []float64{0.1, 0.2, 0.3, 0.4, 0.5},
			context:      "",
			lastExplored: make([]time.Time, 0),
			desiredSpeed: 0.25},
			args{0.55}, 0, false},
		{"findLeftMost_2", fields{
			bins:         []float64{0.1, 0.2, 0.3, 0.4, 0.5},
			context:      "",
			lastExplored: make([]time.Time, 0),
			desiredSpeed: 0.25},
			args{0.05}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uf := UniformFlat{
				bins:         tt.fields.bins,
				context:      tt.fields.context,
				lastExplored: tt.fields.lastExplored,
				desiredSpeed: tt.fields.desiredSpeed,
			}
			got, got1 := uf.findLeftmost(tt.args.price)
			assert.Equalf(t, tt.want, got, "findLeftmost(%v)", tt.args.price)
			assert.Equalf(t, tt.want1, got1, "findLeftmost(%v)", tt.args.price)
		})
	}
}

func TestUniformFlat_sampleBin(t *testing.T) {
	ea := InitUniformFlat("contextHash", 0.12, 3.79, 40, 1, zerolog.New(os.Stdout))
	noOKN := 0
	errN := 0
	sampledBinsN := make([]int, 15)
	for i := 0; i < 300; i++ {
		time.Sleep(5 * time.Millisecond)
		bin, OK, err := ea.sampleBin(1, 15)
		if err != nil {
			errN += 1
			continue
		}
		if !OK {
			noOKN += 1
			continue
		}
		sampledBinsN[bin] += 1
	}
	fmt.Println("noOK: ", noOKN)
	fmt.Println("err: ", errN)
	fmt.Println(sampledBinsN)
	// TODO add empiric assert
}

func TestUniformFlat_sampleNewPrice(t *testing.T) {
	type fields struct {
		bins         []float64
		context      string
		lastExplored []time.Time
		desiredSpeed float64
	}
	type args struct {
		bin int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{

		{"sampleNewPrice_1", fields{[]float64{0.1, 0.2, 0.3, 0.4, 0.5}, "ctx", make([]time.Time, 0), 4.0}, args{0}},
		{"sampleNewPrice_2", fields{[]float64{0.1, 0.2, 0.3, 0.4, 0.5}, "ctx", make([]time.Time, 0), 4.0}, args{1}},
		{"sampleNewPrice_3", fields{[]float64{0.1, 0.2, 0.3, 0.4, 0.5}, "ctx", make([]time.Time, 0), 4.0}, args{2}},
		{"sampleNewPrice_4", fields{[]float64{0.1, 0.2, 0.3, 0.4, 0.5}, "ctx", make([]time.Time, 0), 4.0}, args{3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uf := UniformFlat{
				bins:         tt.fields.bins,
				context:      tt.fields.context,
				lastExplored: tt.fields.lastExplored,
				desiredSpeed: tt.fields.desiredSpeed,
			}
			for i := 0; i < 10; i++ {
				assert.True(t, uf.bins[tt.args.bin] <= uf.sampleNewPrice(tt.args.bin) && uf.bins[tt.args.bin+1] >= uf.sampleNewPrice(tt.args.bin), "sampleNewPrice(%v)", tt.args.bin)
			}
		})
	}
}
