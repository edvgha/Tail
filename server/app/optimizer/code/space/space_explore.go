package space

import "github.com/rs/zerolog/log"

func (s *Space) Explore(floorPrice, price float64) (float64, ExploreData, bool, error) {
	explorationPrice, OK, err := s.explorationAlgorithm.Call(floorPrice, price)
	if err != nil {
		return 0.0, ExploreData{}, false, err
	}
	if !OK {
		return 0.0, ExploreData{}, false, nil
	}

	buckets := s.sampleBuckets(explorationPrice)
	return explorationPrice, ExploreData{s.ContextHash, buckets}, true, nil
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
	log.Debug().Msgf("update: ctx: %s imp: %t", data.ContextHash, impression)

	for i := 0; i < len(data.Buckets); i++ {
		bID := data.Buckets[i]
		if bID == -1 {
			continue
		}
		s.Levels[i].Buckets[bID].Update(impression)
	}
}
