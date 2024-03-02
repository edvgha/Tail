package space

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"math/rand"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"tail.server/app/optimizer/code/misc"
)

const lambdaMin float64 = 0.1
const lambdaMax float64 = 1.8
const minLevelSize int = 1
const minBucketSize int = 5
const minBufferSize int = 10

type Space struct {
	ContextHash          string
	Levels               []*Level
	IsFull               bool
	mutex                sync.Mutex
	explorationAlgorithm Algorithm
	wcMutex              sync.Mutex
	ExplorationQty       atomic.Uint32
	LastUpdateQty        atomic.Uint32
}

type Level struct {
	Buckets      []*Bucket
	WinningCurve []float64
}

func (l *Level) exploit(floorPrice, price float64) (float64, error) {
	left := -1
	right := -1
	for i := 0; i < len(l.Buckets); i++ {
		if floorPrice < l.Buckets[i].Rhs && floorPrice > l.Buckets[i].Lhs {
			left = i
		}
		if price < l.Buckets[i].Rhs && price > l.Buckets[i].Lhs {
			right = i
		}
	}
	if left == -1 || right == -1 || left > right {
		return price, misc.UnfeasiblePriceError{Price: price}
	}
	mid := l.Buckets[left].Lhs + (l.Buckets[left].Rhs-l.Buckets[left].Lhs)/2.0
	if floorPrice > mid {
		left -= 1
	}
	mid = l.Buckets[right].Lhs + (l.Buckets[right].Rhs-l.Buckets[right].Lhs)/2.0
	if price < mid {
		right += 1
	}
	if left < 0 || right >= len(l.Buckets) {
		return price, misc.UnfeasiblePriceError{Price: price}
	}

	maxSoFar := 0.0
	recommendationPrice := 0.0
	for i := left; i <= right; i++ {
		midPrice := l.Buckets[i].Lhs + (l.Buckets[i].Rhs-l.Buckets[i].Lhs)/2.0
		if maxSoFar < (price-midPrice)*l.WinningCurve[i] {
			maxSoFar = (price - midPrice) * l.WinningCurve[i]
			recommendationPrice = midPrice
		}
	}
	return recommendationPrice, nil
}

func (l *Level) sampleBuckets(price float64) int {
	for i := 0; i < len(l.Buckets); i++ {
		if price >= l.Buckets[i].Lhs && price <= l.Buckets[i].Rhs {
			return i
		}
	}
	return -1
}

type Bucket struct {
	Lhs       float64 `json:"lhs"`
	Rhs       float64 `json:"rhs"`
	Alpha     float64 `json:"alpha"`
	Beta      float64 `json:"beta"`
	Buffer    []bool  `json:"buffer"`
	Size      int     `json:"size"`
	Discount  float64 `json:"discount"`
	Pr        float64 `json:"pr"`
	UpdateQty int     `json:"update_qty"`
}

func (b *Bucket) Update(impression bool) {
	b.UpdateQty += 1
	if len(b.Buffer) >= b.Size {
		if b.Buffer[b.Size-1] {
			b.Alpha -= 1
		} else {
			b.Beta -= 1
		}
		buffer := b.Buffer[:b.Size-1]
		b.Buffer = append([]bool{impression}, buffer...)
	} else {
		b.Buffer = append([]bool{impression}, b.Buffer...)
	}

	if impression {
		b.Alpha += 1
	} else {
		b.Beta += 1
	}

	b.Pr = b.Discount*b.Pr + (1-b.Discount)*(b.Alpha/(b.Alpha+b.Beta))
}

type SpaceDesc struct {
	ContextHash string  `json:"context_hash"`
	MinPrice    float64 `json:"min_price"`
	MaxPrice    float64 `json:"max_price"`
}

func NewSpace(contextHash string, minPrice, maxPrice float64, cfg misc.Config) (*Space, error) {
	if minPrice > maxPrice {
		return nil, misc.PriceRangeError{Min: minPrice, Max: maxPrice}
	}

	if cfg.LevelSize < minLevelSize {
		return nil, misc.InvalidLevelError{Level: cfg.LevelSize}
	}

	if cfg.BucketSize < minBucketSize {
		return nil, misc.BucketSizeError{BucketSize: cfg.BucketSize}
	}

	if cfg.BufferSize < minBufferSize {
		return nil, misc.BufferSizeError{BufferSize: cfg.BufferSize}
	}

	if cfg.Discount < 0.0 || cfg.Discount > 1.0 {
		return nil, misc.DiscountFactorError{Discount: cfg.Discount}
	}

	return &Space{
		ContextHash:          contextHash,
		Levels:               newLevels(minPrice, maxPrice, cfg),
		IsFull:               false,
		mutex:                sync.Mutex{},
		explorationAlgorithm: InitUniformFlat(contextHash, minPrice, maxPrice, 2*cfg.BucketSize, cfg.DesiredExplorationSpeed),
		wcMutex:              sync.Mutex{},
	}, nil
}

func newLevels(minPrice, maxPrice float64, cfg misc.Config) []*Level {
	levels := make([]*Level, cfg.LevelSize)

	// slice of lambda parameters for the exp distribution
	lambdas := linspace(cfg.LevelSize)
	for i := 0; i < cfg.LevelSize; i++ {
		levels[i] = newLevel(lambdas[i], minPrice, maxPrice, cfg)
	}
	return levels
}

func linspace(levelSize int) []float64 {
	lambdas := make([]float64, levelSize)
	lambdas[0] = lambdaMin
	if levelSize == 1 {
		return lambdas
	}
	step := (lambdaMax - lambdaMin) / float64(levelSize-1)

	for i := 1; i < levelSize; i++ {
		lambdas[i] = lambdas[i-1] + step
	}
	return lambdas
}

func newLevel(lambda, minPrice, maxPrice float64, cfg misc.Config) *Level {
	buckets := newBuckets(lambda, minPrice, maxPrice, cfg)
	wc := make([]float64, len(buckets))
	for i := 0; i < len(buckets); i++ {
		wc[i] = 0.5
	}
	return &Level{
		Buckets:      buckets,
		WinningCurve: wc,
	}
}

func newBuckets(lambda, minPrice, maxPrice float64, cfg misc.Config) []*Bucket {
	bounds := generateBucketBounds(lambda, minPrice, maxPrice, cfg.BucketSize)
	buckets := make([]*Bucket, cfg.BucketSize)
	lhs := bounds[:len(bounds)-1]
	rhs := bounds[1:]

	for i := 0; i < cfg.BucketSize; i++ {
		buckets[i] = &Bucket{
			Lhs:       lhs[i],
			Rhs:       rhs[i],
			Alpha:     1,
			Beta:      1,
			Buffer:    make([]bool, 0),
			Size:      cfg.BufferSize,
			Discount:  cfg.Discount,
			Pr:        0.5,
			UpdateQty: 0,
		}
	}
	return buckets
}

func generateBucketBounds(lambda, minPrice, maxPrice float64, nBuckets int) []float64 {
	bounds := make([]float64, nBuckets+1)
	// generate random numbers from Exp(lambda) distribution
	for i := 0; i < (nBuckets + 1); i++ {
		bounds[i] = rand.ExpFloat64() / lambda
	}

	sort.Slice(bounds, func(i, j int) bool { return bounds[i] < bounds[j] })

	// scale generated numbers into [minPrice, maxPrice]
	min := bounds[0]
	max := bounds[len(bounds)-1]
	for i := 0; i < (nBuckets + 1); i++ {
		bounds[i] = ((bounds[i]-min)/(max-min))*(maxPrice-minPrice) + minPrice
	}
	return bounds
}

func LoadSpaces(cfg misc.Config) (map[string]*Space, error) {
	file, err := os.Open(cfg.SpaceDescFile)
	if err != nil {
		log.Error().Msgf("Failed to open file %s", cfg.SpaceDescFile)
		return nil, err
	}
	defer file.Close()

	var spacesDesc []SpaceDesc
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&spacesDesc)
	if err != nil {
		log.Error().Msgf("Failed to decode json file: %s", cfg.SpaceDescFile)
		return nil, err
	}

	spaces := make(map[string]*Space)
	for _, s := range spacesDesc {
		spaces[s.ContextHash], err = NewSpace(s.ContextHash, s.MinPrice, s.MaxPrice, cfg)
		if err != nil {
			return nil, err
		}
	}
	log.Debug().Msgf("Description of spaces loaded successfully")
	return spaces, nil
}
