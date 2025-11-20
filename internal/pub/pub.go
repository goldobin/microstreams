package pub

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"

	"github.com/go-redis/redis/v8"

	"github.com/goldobin/microstreams/internal/rng"
)

type Pub struct {
	RedisConfig redis.Options
}

func (p Pub) Publish(idRange rng.Range) func() {
	needStop := new(uint32)

	go func() {
		counter := 0

		rdb := redis.NewClient(&p.RedisConfig)
		defer func() {
			if err := rdb.Close(); err != nil {
				log.Fatalln(err)
			}
		}()

		for {
			v := atomic.LoadUint32(needStop)

			if v > 0 {
				break
			}
			ctx := context.Background()
			id := idRange.PickRandom()
			ch := fmt.Sprintf("ch%06d", id)
			message := fmt.Sprintf("Message %07d from %s", counter, idRange)
			rdb.Publish(ctx, ch, message)

			counter++
		}
	}()

	return func() {
		atomic.StoreUint32(needStop, 1)
	}
}
