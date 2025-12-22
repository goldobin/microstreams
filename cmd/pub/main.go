package main

import (
	"context"
	"math/rand/v2"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"

	"github.com/goldobin/microstreams/internal/rates"

	"github.com/goldobin/microstreams/internal/pub"
	"github.com/goldobin/microstreams/internal/ranges"
)

func main() {
	var (
		logger      = zerolog.New(os.Stdout).Hook(shardHook{}).With().Timestamp().Logger()
		redisClient = newRedisClient()
	)
	defer func(redisClient *redis.Client) {
		if err := redisClient.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close redis client")
		}
	}(redisClient)

	if err := pingRedis(redisClient); err != nil {
		logger.Fatal().Err(err).Msg("Failed to ping redis")
		return
	}
	logger.Print("Redis ping OK")

	const (
		randSeed1            = 1
		randSeed2            = 2
		channelCount         = 20
		publisherCount       = 10
		durationSeconds      = 10
		messageRatePerSecond = 10
	)
	var (
		randSrc                 = rand.NewPCG(randSeed1, randSeed2)
		randGen                 = rand.New(randSrc)
		messageRatePerPublisher = rates.Rate{
			Count:    messageRatePerSecond * 60 * 60 / publisherCount,
			Duration: time.Hour,
		}
		p = pub.Pub{
			Rand:   randGen,
			Redis:  redisClient,
			Logger: logger.With().Str("component", "publisher").Logger(),
			Rate:   messageRatePerPublisher,
		}
		channelRange   = ranges.Int{Left: 0, Right: channelCount}
		channelShards  = channelRange.Split(publisherCount)
		ctx, cancel    = context.WithCancel(context.Background())
		startWg        sync.WaitGroup
		stopWg         sync.WaitGroup
		publishToShard = func(shard ranges.Int) {
			defer stopWg.Done()
			startWg.Done()
			p.Publish(
				context.WithValue(ctx, shardKey, shard),
				shard,
			)
		}
	)
	defer cancel()
	startWg.Add(publisherCount)
	stopWg.Add(publisherCount)

	for _, v := range channelShards {
		go publishToShard(v)
	}

	go func() {
		defer cancel()
		startWg.Wait()
		logger.Print("All publishers ready, starting timer")
		time.Sleep(durationSeconds * time.Second)
	}()

	stopWg.Wait()
}

func newRedisClient() *redis.Client {
	redisCfg := redis.Options{Addr: "localhost:6379"}
	return redis.NewClient(&redisCfg)
}

func pingRedis(r *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if status := r.Ping(ctx); status.Err() != nil {
		return status.Err()
	}

	return nil
}

type (
	contextKey struct{}
	shardHook  struct{}
)

var shardKey = contextKey{}

func (h shardHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	ctx := e.GetCtx()
	shardIdAny := ctx.Value(shardKey)
	shardId, ok := shardIdAny.(ranges.Int)
	if !ok {
		return
	}

	e.Stringer("shard_id", shardId)
}
