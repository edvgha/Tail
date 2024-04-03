package app

import (
	"hash/fnv"
	"strconv"
	"tail.server/app/optimizer/code/misc"
	"tail.server/app/optimizer/code/space"
)

func optimize(request *Request) Response {
	if !validate(request) {
		TLog.Error().Msgf("validation error: r_id %s", request.ID)
		return Response{
			OptimizedPrice: request.Price,
			Status:         "validation error",
		}
	}
	explorationPrice, OK, err := explore(request)
	if err != nil {
		return Response{
			OptimizedPrice: request.Price,
			Status:         err.Error(),
		}
	}
	if !OK {
		recommendedPrice, err := exploit(request)
		if err != nil {
			return Response{
				OptimizedPrice: request.Price,
				Status:         err.Error(),
			}
		}
		return Response{
			OptimizedPrice: recommendedPrice,
			Status:         "exploited",
		}
	}
	return Response{
		OptimizedPrice: explorationPrice,
		Status:         "explored",
	}
}

func explore(request *Request) (float64, bool, error) {
	context := contextHash(request)
	s, exists := Spaces[context]
	if !exists {
		return 0.0, false, misc.NoSpaceError{}
	}

	newPrice, data, OK, err := s.Explore(request.FloorPrice, request.Price)
	if err != nil {
		TLog.Error().Msgf("explore ctx: %s error: %s", context, err.Error())
		return 0.0, false, err
	}
	s.ExplorationQty.Add(1)
	if !OK {
		TLog.Trace().Msgf("ctx: %s NO exploration for price: %f ", context, request.Price)
		return 0.0, false, nil
	}
	data.ContextHash = context
	cb := func(d space.ExploreData) bool {
		s, exists := Spaces[data.ContextHash]
		if !exists {
			return false
		}
		s.Update(data, false)
		return true
	}
	Cache.Set(request.ID, data, CacheTTL, cb)
	TLog.Trace().Msgf("explore: %s price: %f new_price: %f", context, request.Price, newPrice)
	return newPrice, true, nil
}

func exploit(request *Request) (float64, error) {
	context := contextHash(request)
	space, exists := Spaces[context]
	if !exists {
		return 0.0, misc.NoSpaceError{}
	}

	recommendedPrice, err := space.Exploit(request.FloorPrice, request.Price)
	if err != nil {
		TLog.Error().Msgf("exploit ctx: %s error: %s", context, err.Error())
		return request.Price, err
	}
	return recommendedPrice, nil
}

func validate(request *Request) bool {
	if request.FloorPrice <= 0 {
		return false
	}
	if request.Price <= 0 {
		return false
	}
	if request.FloorPrice > request.Price {
		return false
	}
	return true
}

func contextHash(request *Request) string {
	h := fnv.New64a()
	h.Write([]byte(request.DC))
	h.Write([]byte(request.BundleID))
	h.Write([]byte(request.TagID))
	h.Write([]byte(request.GeoCountry))
	h.Write([]byte(request.AdFormat))
	h.Write([]byte(request.AdFormat))
	return strconv.FormatUint(h.Sum64(), 10)
}
