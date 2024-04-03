package space

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"tail.server/app/optimizer/code/misc"
	"testing"
)

func Test_linspace(t *testing.T) {
	tests := []struct {
		name string
		args int
		want []float64
	}{
		{"lambda_1", 1, []float64{lambdaMin}},
		{"lambda_2", 2, []float64{lambdaMin, lambdaMax}},
		{"lambda_3", 3, []float64{lambdaMin, 0.95, 1.7999999999999998}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := linspace(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("linspace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateBucketBounds(t *testing.T) {
	type args struct {
		lambda   float64
		minPrice float64
		maxPrice float64
		nBuckets int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"BB_0", args{lambdaMax, 0.1, 1, 2}, 3},
		{"BB_1", args{lambdaMax, 0.1, 1, 10}, 11},
		{"BB_2", args{lambdaMin, 0.1, 1, 10}, 11},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateBucketBounds(tt.args.lambda, tt.args.minPrice, tt.args.maxPrice, tt.args.nBuckets)
			if len(got) != tt.want {
				t.Errorf("generateBucketBounds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newBuckets(t *testing.T) {
	type args struct {
		lambda   float64
		minPrice float64
		maxPrice float64
		cfg      misc.Config
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"newBuckets_0", args{lambdaMin, 0.1, 10.0, misc.Config{BucketSize: 5, BufferSize: 10, Discount: 0.2}}, 5},
		{"newBuckets_1", args{lambdaMax, 0.1, 21.0, misc.Config{BucketSize: 50, BufferSize: 100, Discount: 0.25}}, 50},
		{"newBuckets_2", args{lambdaMin, 0.05, 33.0, misc.Config{BucketSize: 100, BufferSize: 100, Discount: 0.75}}, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newBuckets(tt.args.lambda, tt.args.minPrice, tt.args.maxPrice, tt.args.cfg)
			if len(got) != tt.want {
				t.Errorf("newBuckets() = %v, want %v", got, tt.want)
			}
			if got[0].Discount != tt.args.cfg.Discount {
				t.Errorf("NewBuckets() Discount = %v want %v", got[0].Discount, tt.args.cfg.Discount)
			}
			if got[0].Size != tt.args.cfg.BufferSize {
				t.Errorf("NewBuckets() BufferSize = %v want %v", got[0].Size, tt.args.cfg.BufferSize)
			}
			if got[0].Pr != 0.5 {
				t.Errorf("NewBuckets() Pr = %v want 0.5", got[0].Pr)
			}
			if got[0].Lhs != tt.args.minPrice {
				t.Errorf("NewBuckets() Lhs[0] = %v want %v", got[0].Lhs, tt.args.minPrice)
			}
			if got[len(got)-1].Rhs != tt.args.maxPrice {
				t.Errorf("NewBuckets() Rhs[len(BucketSize) - 1] = %v want %v", got[len(got)-1].Rhs, tt.args.maxPrice)
			}
			for i := 0; i < len(got); i++ {
				if got[i].Lhs > got[i].Rhs {
					t.Errorf("NewBuckets() invalid bucket: index = %d Lhs = %v Rhs = %v", i, got[i].Lhs, got[i].Rhs)
				}
				if i > 0 {
					if got[i-1].Rhs != got[i].Lhs {
						t.Errorf("NewBuckets() invalid neighbor buckets index left = %d, index right %d", (i - 1), i)
					}
				}
			}
		})
	}
}

func Test_newLevels(t *testing.T) {
	type args struct {
		minPrice float64
		maxPrice float64
		cfg      misc.Config
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"newLevels_0", args{0.2, 31.4, misc.Config{LevelSize: 10, BucketSize: 100, BufferSize: 100, Discount: 0.75}}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newLevels(tt.args.minPrice, tt.args.maxPrice, tt.args.cfg); len(got) != tt.want {
				t.Errorf("newLevels() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSpace(t *testing.T) {
	type args struct {
		contextHash string
		minPrice    float64
		maxPrice    float64
		cfg         misc.Config
	}
	tests := []struct {
		name string
		args args
	}{
		{"NewSpace_0", args{"871397612603656680", 0.2, 22.1, misc.Config{LevelSize: 10, BucketSize: 100, BufferSize: 100, Discount: 0.75}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSpace(tt.args.contextHash, tt.args.minPrice, tt.args.maxPrice, tt.args.cfg, zerolog.New(os.Stdout))
			if err != nil {
				t.Errorf("%s", err.Error())
			}
			if got.ContextHash != tt.args.contextHash || got.IsFull == true || len(got.Levels) != tt.args.cfg.LevelSize {
				t.Errorf("Invalid space")
			}
		})
	}
}

func TestLoadSpaces(t *testing.T) {
	cfg := misc.Config{SpaceDescFile: "spaces_desc.json", LevelSize: 10, BucketSize: 100, BufferSize: 100, Discount: 0.75}
	spaces, err := LoadSpaces(cfg, zerolog.New(os.Stdout))
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	assert.Equal(t, len(spaces), 7, "expected number of spaces is 7")
	space := spaces["37811cc60c99ec9a"]
	assert.NotNilf(t, space, "expected to not be nil")
	assert.Equal(t, len(space.Levels), 10, "expected number of levels is 10")
	assert.Equal(t, len(space.Levels[0].Buckets), 100, "expected number of buckets is 100")
	assert.Equal(t, space.Levels[0].Buckets[0].Lhs, 0.2035, "")
	assert.Equal(t, space.Levels[9].Buckets[99].Rhs, 0.4485, "")
}

func TestLevel_sampleBuckets(t *testing.T) {
	type fields struct {
		Buckets []*Bucket
	}
	type args struct {
		price float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{"sampleBuckets_1",
			fields{
				Buckets: newBuckets(0.1, 0.2, 9.8, misc.Config{BucketSize: 100, BufferSize: 10, Discount: 0.2}),
			},
			args{
				price: 0.3,
			},
			10},
		{"sampleBuckets_2",
			fields{
				Buckets: newBuckets(0.1, 0.2, 9.8, misc.Config{BucketSize: 100, BufferSize: 10, Discount: 0.2}),
			},
			args{
				price: 0.21,
			},
			10},
		{"sampleBuckets_3",
			fields{
				Buckets: newBuckets(0.1, 0.2, 9.8, misc.Config{BucketSize: 100, BufferSize: 10, Discount: 0.2}),
			},
			args{
				price: 0.27,
			},
			10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Level{
				Buckets: tt.fields.Buckets,
			}
			assert.Greaterf(t, tt.want, l.sampleBuckets(tt.args.price), "sampleBuckets(%v)", tt.args.price)
		})
	}
}

func TestBucket_Update(t *testing.T) {
	b := &Bucket{
		Lhs:       0.1,
		Rhs:       0.2,
		Alpha:     1,
		Beta:      1,
		Buffer:    make([]bool, 0),
		Size:      5,
		Discount:  0.25,
		Pr:        0.5,
		UpdateQty: 0,
	}
	type args struct {
		impression bool
	}
	tests := []struct {
		name string
		args args
	}{
		{"Update", args{impression: true}},
		{"Update", args{impression: true}},
		{"Update", args{impression: true}},
		{"Update", args{impression: true}},
		{"Update", args{impression: true}},
		{"Update", args{impression: true}},
		{"Update", args{impression: true}},
		{"Update", args{impression: true}},
	}
	for i, tt := range tests {
		i := i
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			b.Update(tt.args.impression)
		})
		assert.Equal(t, i+1, b.UpdateQty)
		if i >= 5 {
			assert.Equal(t, 5, len(b.Buffer))
		} else {
			assert.Equal(t, i+1, len(b.Buffer))
		}
	}
}

func TestLevel_exploit(t *testing.T) {
	buckets := make([]*Bucket, 10)
	for i, v := range []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
		buckets[i] = &Bucket{Lhs: v, Rhs: v + 1}
	}
	wc := []float64{0.05, 0.15, 0.25, 0.35, 0.45, 0.55, 0.65, 0.75, 0.85, 0.95}
	type fields struct {
		Buckets      []*Bucket
		WinningCurve []float64
	}
	type args struct {
		floorPrice float64
		price      float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{"exploit_0", fields{Buckets: buckets, WinningCurve: wc}, args{1.2, 9.8}, 5.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Level{
				Buckets:      tt.fields.Buckets,
				WinningCurve: tt.fields.WinningCurve,
			}
			got, err := l.exploit(tt.args.floorPrice, tt.args.price, zerolog.New(os.Stdout))
			assert.Nil(t, err)
			assert.Equalf(t, tt.want, got, "exploit(%v, %v)", tt.args.floorPrice, tt.args.price)
		})
	}
}

func TestSpace_WC(t *testing.T) {
	s, err := NewSpace("ctxHash", 0.1, 9.9, misc.Config{
		BufferSize:              100,
		Discount:                0.25,
		DesiredExplorationSpeed: 0.2,
		LogLevel:                "debug",
		LevelSize:               3,
		BucketSize:              10,
		SpaceDescFile:           "file.txt",
		CacheTTL:                1,
	}, zerolog.New(os.Stdout))
	assert.Nil(t, err)
	le := s.WC()
	assert.Equal(t, 3, len(le.Level))
	assert.True(t, reflect.DeepEqual(le.Level[0].Pr, s.Levels[0].WinningCurve))
	assert.True(t, reflect.DeepEqual(le.Level[1].Pr, s.Levels[1].WinningCurve))
	assert.True(t, reflect.DeepEqual(le.Level[2].Pr, s.Levels[2].WinningCurve))
	// -----------------------------------------------------------------------
	prices := make([]float64, len(s.Levels[0].Buckets))
	for i := 0; i < len(prices); i++ {
		prices[i] = s.Levels[0].Buckets[i].Lhs + (s.Levels[0].Buckets[i].Rhs-s.Levels[0].Buckets[i].Lhs)/2.0
	}
	assert.True(t, reflect.DeepEqual(le.Level[0].Price, prices))
	for i := 0; i < len(prices); i++ {
		prices[i] = s.Levels[1].Buckets[i].Lhs + (s.Levels[1].Buckets[i].Rhs-s.Levels[1].Buckets[i].Lhs)/2.0
	}
	assert.True(t, reflect.DeepEqual(le.Level[1].Price, prices))
	for i := 0; i < len(prices); i++ {
		prices[i] = s.Levels[2].Buckets[i].Lhs + (s.Levels[2].Buckets[i].Rhs-s.Levels[2].Buckets[i].Lhs)/2.0
	}
	assert.True(t, reflect.DeepEqual(le.Level[2].Price, prices))
}
