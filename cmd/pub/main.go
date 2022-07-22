package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"gitlab.newmotion.com/cpo/incubator/microstreams/internal/pub"
	"gitlab.newmotion.com/cpo/incubator/microstreams/internal/rng"
	"time"
)

func main() {
	fmt.Println("Pub")

	redisConfig := redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}

	p := pub.Pub{
		RedisConfig: redisConfig,
	}

	rng := rng.Range{
		Left:  0,
		Right: 140_000,
	}

	ranges := rng.Split(10)
	stops := make([]func(), len(ranges))

	for i, r := range ranges {
		stops[i] = p.Publish(r)
	}

	defer func() {
		for _, stop := range stops {
			stop()
		}
	}()

	time.Sleep(time.Minute * 1)
}
