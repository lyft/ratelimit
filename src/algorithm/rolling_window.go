package algorithm

import (
	"math"

	"github.com/coocood/freecache"
	pb "github.com/envoyproxy/go-control-plane/envoy/service/ratelimit/v3"
	"github.com/envoyproxy/ratelimit/src/config"
	"github.com/envoyproxy/ratelimit/src/utils"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/durationpb"
)

const DummyCacheKeyTime = 0

type RollingWindowImpl struct {
	timeSource        utils.TimeSource
	cacheKeyGenerator utils.CacheKeyGenerator
	localCache        *freecache.Cache
	nearLimitRatio    float32
	arrivedAt         int64
	newTat            int64
}

func (rw *RollingWindowImpl) GetResponseDescriptorStatus(key string, limit *config.RateLimit, results int64, isOverLimitWithLocalCache bool, hitsAddend int64) *pb.RateLimitResponse_DescriptorStatus {
	if key == "" {
		return &pb.RateLimitResponse_DescriptorStatus{
			Code:           pb.RateLimitResponse_OK,
			CurrentLimit:   nil,
			LimitRemaining: 0,
		}
	}
	if isOverLimitWithLocalCache {
		rw.PopulateStats(limit, 0, uint64(hitsAddend), uint64(hitsAddend))
		return &pb.RateLimitResponse_DescriptorStatus{
			Code:               pb.RateLimitResponse_OVER_LIMIT,
			CurrentLimit:       limit.Limit,
			LimitRemaining:     0,
			DurationUntilReset: utils.CalculateFixedReset(limit.Limit, rw.timeSource),
		}
	}

	isOverLimit, limitRemaining, durationUntilReset := rw.IsOverLimit(limit, int64(results), hitsAddend)
	if !isOverLimit {
		return &pb.RateLimitResponse_DescriptorStatus{
			Code:               pb.RateLimitResponse_OK,
			CurrentLimit:       limit.Limit,
			LimitRemaining:     uint32(limitRemaining),
			DurationUntilReset: durationUntilReset,
		}
	} else {
		if rw.localCache != nil {
			err := rw.localCache.Set([]byte(key), []byte{}, int(utils.UnitToDivider(limit.Limit.Unit)))
			if err != nil {
				logger.Errorf("Failing to set local cache key: %s", key)
			}
		}

		return &pb.RateLimitResponse_DescriptorStatus{
			Code:               pb.RateLimitResponse_OVER_LIMIT,
			CurrentLimit:       limit.Limit,
			LimitRemaining:     uint32(limitRemaining),
			DurationUntilReset: durationUntilReset,
		}
	}
}

func (rw *RollingWindowImpl) GetNewTat() int64 {
	return rw.newTat
}
func (rw *RollingWindowImpl) GetArrivedAt() int64 {
	return rw.newTat
}

func (rw *RollingWindowImpl) IsOverLimit(limit *config.RateLimit, results int64, hitsAddend int64) (bool, int64, *durationpb.Duration) {
	now := rw.timeSource.UnixNanoNow()

	// Time during computation should be in nanosecond
	rw.arrivedAt = now
	tat := utils.MaxInt64(results, rw.arrivedAt)
	totalLimit := int64(limit.Limit.RequestsPerUnit)
	period := utils.SecondsToNanoseconds(utils.UnitToDivider(limit.Limit.Unit))
	quantity := int64(hitsAddend)

	// GCRA computation
	// Emission interval is the cost of each request
	emissionInterval := period / totalLimit
	// Tat is set to current request timestamp if not set before

	// New tat define the end of the window
	rw.newTat = tat + emissionInterval*quantity
	// We allow the request if it's inside the window
	allowAt := rw.newTat - period
	diff := rw.arrivedAt - allowAt

	previousAllowAt := tat - period
	previousLimitRemaining := int64(math.Ceil(float64((rw.arrivedAt - previousAllowAt) / emissionInterval)))
	previousLimitRemaining = utils.MaxInt64(previousLimitRemaining, 0)
	nearLimitWindow := int64(math.Ceil(float64(float32(limit.Limit.RequestsPerUnit) * (1.0 - rw.nearLimitRatio))))
	limitRemaining := int64(math.Ceil(float64(diff / emissionInterval)))
	hitNearLimit := quantity - (utils.MaxInt64(previousLimitRemaining, nearLimitWindow) - nearLimitWindow)

	if diff < 0 {
		rw.PopulateStats(limit, uint64(utils.MinInt64(previousLimitRemaining, nearLimitWindow)), uint64(quantity-previousLimitRemaining), 0)

		return true, 0, utils.NanosecondsToDuration(int64(math.Ceil(float64(tat - rw.arrivedAt))))
	} else {
		if hitNearLimit > 0 {
			rw.PopulateStats(limit, uint64(hitNearLimit), 0, 0)
		}

		return false, limitRemaining, utils.NanosecondsToDuration(rw.newTat - rw.arrivedAt)
	}
}

func (rw *RollingWindowImpl) IsOverLimitWithLocalCache(key string) bool {
	if rw.localCache != nil {
		_, err := rw.localCache.Get([]byte(key))
		if err == nil {
			return true
		}
	}
	return false
}

func (rw *RollingWindowImpl) GenerateCacheKeys(request *pb.RateLimitRequest,
	limits []*config.RateLimit, hitsAddend int64) []utils.CacheKey {
	return rw.cacheKeyGenerator.GenerateCacheKeys(request, limits, uint32(hitsAddend), DummyCacheKeyTime)
}

func (rw *RollingWindowImpl) PopulateStats(limit *config.RateLimit, nearLimit uint64, overLimit uint64, overLimitWithLocalCache uint64) {
	limit.Stats.NearLimit.Add(nearLimit)
	limit.Stats.OverLimit.Add(overLimit)
	limit.Stats.OverLimitWithLocalCache.Add(overLimitWithLocalCache)
}

func NewRollingWindowAlgorithm(timeSource utils.TimeSource, localCache *freecache.Cache, nearLimitRatio float32) *RollingWindowImpl {
	return &RollingWindowImpl{
		timeSource:        timeSource,
		cacheKeyGenerator: utils.NewCacheKeyGenerator(),
		localCache:        localCache,
		nearLimitRatio:    nearLimitRatio,
	}
}
