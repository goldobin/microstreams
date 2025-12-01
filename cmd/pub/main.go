package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"github.com/goldobin/microstreams/internal/rates"

	"github.com/goldobin/microstreams/internal/pub"
	"github.com/goldobin/microstreams/internal/ranges"
)

func main() {
	var (
		log, logSync = newLogger()
		redisCfg     = redis.Options{Addr: "localhost:6379"}
		rdb          = redis.NewClient(&redisCfg)
	)
	defer logSync()

	if err := pingRedis(rdb); err != nil {
		log.Fatal("Failed to ping redis", zap.Error(err))
	}
	log.Info("Redis ping OK")

	const (
		randSeed1       = 1
		randSeed2       = 2
		channelCount    = 20
		shardCount      = 10
		durationSeconds = 10
	)
	var (
		randSrc        = rand.NewPCG(randSeed1, randSeed2)
		randGen        = rand.New(randSrc)
		p              = pub.Pub{Rand: randGen, Redis: rdb, Log: log.Named("publisher")}
		channelRange   = ranges.IntRange{Left: 0, Right: channelCount}
		messageRate    = rates.New(100, time.Second)
		ctx, cancel    = context.WithCancel(context.Background())
		channelShards  = channelRange.Split(shardCount)
		startWg        sync.WaitGroup
		stopWg         sync.WaitGroup
		publishToShard = func(rg ranges.IntRange) {
			defer stopWg.Done()
			startWg.Done()
			p.Publish(ctx, rg, messageRate)
		}
	)
	defer cancel()
	startWg.Add(shardCount)
	stopWg.Add(shardCount)

	for _, v := range channelShards {
		publishToShard(v)
	}

	go func() {
		defer cancel()
		startWg.Wait()
		log.Info("All publishers ready, starting timer")
		time.Sleep(durationSeconds * time.Second)
	}()

	stopWg.Wait()
}

func pingRedis(r *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if status := r.Ping(ctx); status.Err() != nil {
		return status.Err()
	}

	return nil
}

func newLogger() (*zap.Logger, func()) {
	var (
		log     = zap.Must(zap.NewProduction())
		logSync = func() {
			if err := log.Sync(); err != nil {
				fmt.Printf("Failed to logSync logger: %v", err)
			}
		}
	)

	return log, logSync
}
