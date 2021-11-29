package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/eflem00/go-example-app/gateways/cache"
	"github.com/eflem00/go-example-app/gateways/db"
	"github.com/rs/zerolog/log"
)

// check cache for key and touch if we get a cache hit
// if cache miss, go to persistant storage and set
func GetResultById(ctx context.Context, key string) (string, error) {
	cacheClient := cache.NewClient()

	val, err := cacheClient.Get(ctx, key).Result()

	// should check the type of error for redis.Nil here but we'll keep it simple and treat this as a cache miss
	if err != nil {
		log.Debug().Msgf("Cache miss for key %v", key)

		val, err = db.GetResultById(key)

		if err != nil {
			return "", errors.New("no value for provided key")
		}

		cacheClient.Set(ctx, key, val, time.Second*5).Result()

		return val, nil
	} else { // cache hit, use the value and touch the key
		log.Debug().Msgf("Cache hit for key %v", key)

		cacheClient.Touch(ctx, key).Err()

		return val, nil
	}
}

func WriteResult(ctx context.Context, key string, value string) error {
	return db.WriteResult(key, value)
}
