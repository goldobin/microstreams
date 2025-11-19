package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/goldobin/microstreams/internal/pub"
	"github.com/goldobin/microstreams/internal/rng"
)

func main() {
	fmt.Println("Pub")

	var (
		redisConfig = redis.Options{Addr: "localhost:6379"}
		p           = pub.Pub{RedisConfig: redisConfig}
		r           = _range.Range{Left: 0, Right: 140_000}
		ranges      = r.Split(10)
		stops       = make([]func(), len(ranges))
	)
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
