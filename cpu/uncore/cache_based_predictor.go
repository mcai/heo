package uncore

import (
	"github.com/mcai/heo/cpu/mem"
	"github.com/mcai/heo/simutil"
)

type CacheBasedPredictorLineValueProvider struct {
	*BaseCacheLineStateProvider
	PredictedValue interface{}
	Confidence     *simutil.SaturatingCounter
}

func NewCacheBasedPredictorLineValueProvider(counterThreshold uint32, counterMaxValue uint32) *CacheBasedPredictorLineValueProvider {
	var valueProvider = &CacheBasedPredictorLineValueProvider{
		BaseCacheLineStateProvider: NewBaseCacheLineStateProvider(
			false,
			func(state interface{}) bool {
				return state != false
			},
		),
		Confidence: simutil.NewSaturatingCounter(0, counterThreshold, counterMaxValue, 0),
	}

	return valueProvider
}

type CacheBasedPredictor struct {
	Cache            *EvictableCache
	NumHits          uint32
	NumMisses        uint32
	NumFailedPredict uint32
}

func NewCacheBasedPredictor(capacity uint32, counterThreshold uint32, counterMaxValue uint32) *CacheBasedPredictor {
	var predictor = &CacheBasedPredictor{
		Cache: NewEvictableCache(
			mem.NewGeometry(
				capacity,
				capacity,
				1,
			),
			func(set uint32, way uint32) CacheLineStateProvider {
				return NewCacheBasedPredictorLineValueProvider(counterThreshold, counterMaxValue)
			},
			CacheReplacementPolicyType_LRU,
		),
	}

	return predictor
}

func (predictor *CacheBasedPredictor) Predict(address uint32, defaultValue interface{}) interface{} {
	var lineFound = predictor.Cache.FindLine(address)

	if lineFound != nil {
		var stateProvider = lineFound.StateProvider.(*CacheBasedPredictorLineValueProvider)

		if stateProvider.Confidence.Taken() {
			return stateProvider.PredictedValue
		}
	}

	predictor.NumFailedPredict++
	return defaultValue
}

func (predictor *CacheBasedPredictor) Update(address uint32, observedValue interface{}) {
	if predictor.Predict(address, nil) == observedValue {
		predictor.NumHits++
	} else {
		predictor.NumMisses++
	}

	var set = predictor.Cache.GetSet(address)
	var tag = predictor.Cache.GetTag(address)

	var cacheAccess = predictor.Cache.NewAccess(nil, address)

	var line = predictor.Cache.Sets[set].Lines[cacheAccess.Way]
	var stateProvider = line.StateProvider.(*CacheBasedPredictorLineValueProvider)

	if cacheAccess.HitInCache {
		if stateProvider.PredictedValue == observedValue {
			stateProvider.Confidence.Update(true)
		} else {
			if stateProvider.Confidence.Value() == 0 {
				stateProvider.PredictedValue = observedValue
			} else {
				stateProvider.Confidence.Update(false)
			}
		}

		predictor.Cache.ReplacementPolicy.HandlePromotionOnHit(nil, set, cacheAccess.Way)
	} else {
		stateProvider.SetState(true)
		line.Tag = int32(tag)

		stateProvider.PredictedValue = observedValue
		stateProvider.Confidence.Reset()

		predictor.Cache.ReplacementPolicy.HandleInsertionOnMiss(nil, set, cacheAccess.Way)
	}
}
