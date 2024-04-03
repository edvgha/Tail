package space

import (
	"github.com/rs/zerolog"
	"math/rand"
	"sort"
	"sync"
	"tail.server/app/optimizer/code/misc"
	"time"
)

type ExploreData struct {
	ContextHash string
	Buckets     []int
	started     time.Time
}

type Algorithm interface {
	Call(floorPrice, price float64) (float64, bool, error)
}

type UniformFlat struct {
	bins         []float64
	context      string
	lastExplored []time.Time
	desiredSpeed float64
	mutex        *sync.Mutex
	log          zerolog.Logger
}

func InitUniformFlat(context string, minPrice, maxPrice float64, nBins int, desiredSpeed float64, log zerolog.Logger) UniformFlat {
	bins := make([]float64, nBins+1)
	lastExplored := make([]time.Time, nBins)
	bins[0] = minPrice
	lastExplored[0] = time.Now()
	step := (maxPrice - minPrice) / float64(nBins)
	for i := 1; i <= nBins; i++ {
		bins[i] = bins[i-1] + step
		if i != nBins {
			lastExplored[i] = time.Now()
		}
	}

	return UniformFlat{
		bins:         bins,
		context:      context,
		lastExplored: lastExplored,
		desiredSpeed: desiredSpeed,
		mutex:        &sync.Mutex{},
		log:          log,
	}
}

func (uf UniformFlat) Call(floorPrice, price float64) (float64, bool, error) {
	r, OK := uf.findLeftmost(price)
	if !OK {
		return 0.0, false, misc.UnfeasiblePriceError{Price: price, Min: uf.bins[0], Max: uf.bins[len(uf.bins)-1]}
	}

	l, OK := uf.findLeftmost(floorPrice)
	if !OK {
		return 0.0, false, misc.UnfeasiblePriceError{Price: floorPrice, Min: uf.bins[0], Max: uf.bins[len(uf.bins)-1]}
	}

	if r-l < 2 {
		return 0.0, false, nil
	}
	bin, OK, err := uf.sampleBin(l, r)
	if err != nil {
		return 0.0, false, err
	}
	if !OK {
		return 0.0, false, nil
	}
	newPrice := uf.sampleNewPrice(bin)
	return newPrice, true, nil
}

func (uf UniformFlat) sampleNewPrice(bin int) float64 {
	l := uf.bins[bin]
	r := uf.bins[bin+1]
	return l + (r-l)*rand.Float64()
}

func (uf UniformFlat) sampleBin(l, r int) (int, bool, error) {
	uf.mutex.Lock()
	defer uf.mutex.Unlock()
	if r-l == 2 {
		t := time.Since(uf.lastExplored[l+1]).Seconds()
		estimatedSpeed := 1 / t
		if rand.Float64() < uf.desiredSpeed/estimatedSpeed {
			uf.lastExplored[l+1] = time.Now()
			return l + 1, true, nil
		} else {
			return l + 1, false, nil
		}
	}
	type elem struct {
		bin   int
		speed float64
	}
	s := make([]elem, r-l-1)
	for i := l + 1; i < r; i++ {
		t := time.Since(uf.lastExplored[i]).Seconds()
		estimatedSpeed := 1 / t
		s[i-l-1] = elem{bin: i, speed: estimatedSpeed}
	}
	sort.Slice(s, func(i, j int) bool {
		return s[i].speed < s[j].speed
	})
	for i := 0; i < len(s); i++ {
		if rand.Float64() < uf.desiredSpeed/s[i].speed {
			uf.lastExplored[s[i].bin] = time.Now()
			return s[i].bin, true, nil
		}
	}
	return 0.0, false, nil
}

func (uf UniformFlat) findLeftmost(price float64) (int, bool) {
	if price <= uf.bins[0] || price > uf.bins[len(uf.bins)-1] {
		uf.log.Error().Msgf("unfeasible price: %v$, for [%v$, %v$]", price, uf.bins[0], uf.bins[len(uf.bins)-1])
		return 0, false
	}
	return sort.SearchFloat64s(uf.bins, price) - 1, true
}
