// This file is manually authored.

package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Charles546/spicker/v2/restapi/operations"
	redis "github.com/go-redis/redis/v8"
)

type CacheStub struct {
	redisClient *redis.Client
}

var redisClient *CacheStub

func getRedisClient() *redis.Client {
	if redisClient == nil {
		redisClient = &CacheStub{}

		connStr := os.Getenv("REDIS_CONNECTION")
		if len(connStr) == 0 {
			log.Printf("no cache setting found in environment")

			return nil
		}

		opts, e := redis.ParseURL(connStr)
		if e != nil {
			log.Printf("invalid redis connection string, no cache setup\n")

			return nil
		}

		redisClient.redisClient = redis.NewClient(opts)
	}

	return redisClient.redisClient
}

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second)
}

func FromCache(symbol string, ndays int) *operations.StockpricesOKBody {
	c := getRedisClient()
	if c == nil {
		return nil
	}

	ctx, cancel := getContext()
	defer cancel()
	b, e := c.Get(ctx, fmt.Sprintf("%s_%d", symbol, ndays)).Result()
	if e != nil {
		if errors.Is(e, redis.Nil) {
			log.Printf("redis key not found: %s_%d\n", symbol, ndays)
		} else {
			log.Printf("server cache error: %+v\n", e)
		}

		return nil
	}

	ret := operations.StockpricesOKBody{}
	e = json.Unmarshal([]byte(b), &ret)
	if e != nil {
		log.Printf("server cache unmarshal error: %+v\n", e)

		return nil
	}

	return &ret
}

func Cache(symbol string, ndays int, payload *operations.StockpricesOKBody) {
	c := getRedisClient()
	if c == nil {
		return
	}

	newYork, _ := time.LoadLocation("America/New_York")
	now := time.Now().In(newYork)
	available := time.Date(now.Year(), now.Month(), now.Day(), 20, 0, 0, 0, newYork)
	if available.Before(now) {
		available.AddDate(0, 0, 1)
	}

	b, e := json.Marshal(payload)
	if e != nil {
		log.Printf("server cache marshal error: %+v\n", e)

		return
	}

	ctx, cancel := getContext()
	defer cancel()
	_, e = c.SetEX(ctx, fmt.Sprintf("%s_%d", symbol, ndays), b, available.Sub(now)).Result()
	if e != nil {
		log.Printf("server cache save error: %+v\n", e)
	}

	return
}
