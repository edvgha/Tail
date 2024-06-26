package space

import (
	"time"
)

func (s *Space) Explore(floorPrice, price float64) (float64, ExploreData, time.Duration, bool, error) {
	t := time.Now()
	explorationPrice, OK, err := s.explorationAlgorithm.Call(floorPrice, price)
	if err != nil {
		return 0.0, ExploreData{}, 0, false, err
	}
	if !OK {
		return 0.0, ExploreData{}, 0, false, nil
	}

	buckets := s.sampleBuckets(explorationPrice)
	return explorationPrice, ExploreData{s.ContextHash, buckets, t}, s.ttl.Time(), true, nil
}

func (s *Space) sampleBuckets(price float64) []int {
	buckets := make([]int, len(s.Levels))
	for i := 0; i < len(s.Levels); i++ {
		buckets[i] = s.Levels[i].sampleBuckets(price)
	}
	return buckets
}

func (s *Space) Update(data ExploreData, impression bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.log.Debug().Msgf("update: ctx: %s imp: %t", data.ContextHash, impression)

	for i := 0; i < len(data.Buckets); i++ {
		bID := data.Buckets[i]
		if bID == -1 {
			continue
		}
		s.Levels[i].Buckets[bID].Update(impression)
	}
	if impression {
		s.log.Debug().Msgf("ack time %v", time.Since(data.started))
		s.ttl.Add(time.Since(data.started))
	}
}
