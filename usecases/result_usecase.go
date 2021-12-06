package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/eflem00/go-example-app/common"
	"github.com/eflem00/go-example-app/gateways/cache"
	"github.com/eflem00/go-example-app/gateways/db"
)

type ResultUsecase struct {
	logger           *common.Logger
	cache            *cache.Cache
	resultRepository *db.ResultRepository
}

func NewResultUseCase(logger *common.Logger, cache *cache.Cache, resultRepository *db.ResultRepository) *ResultUsecase {
	return &ResultUsecase{
		logger,
		cache,
		resultRepository,
	}
}

// check cache for key and touch if we get a cache hit
// if cache miss, go to persistant storage and set
func (uc *ResultUsecase) GetResultById(ctx context.Context, key string) (string, error) {
	val, err := uc.cache.Get(ctx, key)

	// should check the type of error for redis.Nil here but we'll keep it simple and treat this as a cache miss
	if err != nil {
		uc.logger.Debugf("Cache miss for key %v", key)

		val, err = uc.resultRepository.GetResultById(key)

		if err != nil {
			return "", errors.New("no value for provided key")
		}

		uc.cache.Set(ctx, key, val, time.Hour)

		return val, nil
	} else { // cache hit, use the value and touch the key
		uc.logger.Debugf("Cache hit for key %v", key)

		uc.cache.Touch(ctx, key)

		return val, nil
	}
}

func (uc *ResultUsecase) WriteResult(ctx context.Context, key string, value string) error {
	uc.cache.Set(ctx, key, value, time.Hour)

	return uc.resultRepository.WriteResult(key, value)
}
