package pub

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"github.com/goldobin/microstreams/internal/rates"

	"github.com/goldobin/microstreams/internal/ranges"
)

type Pub struct {
	Rand  *rand.Rand
	Rate  rates.Rate
	Redis *redis.Client
	Log   *zap.Logger
}

func (p Pub) Publish(ctx context.Context, r ranges.IntRange, rt rates.Rate) {
	time.Sleep(rt.Random(p.Rand))
	var (
		counter = 0
		log     = p.Log.With(zap.String("range", r.String()))
		ticker  = time.NewTicker(rt.Interval())
	)
	defer ticker.Stop()
	log.Info("Starting publish...")
	defer log.Info("Publish stopped")

	for {
		counter++
		select {
		case <-ctx.Done():
			log.Debug("Stopped", zap.String("stop_reason", ctx.Err().Error()))
			return
		case _, ok := <-ticker.C:
			if !ok {
				log.Debug("Stopped", zap.String("stop_reason", "ticker stopped"))
				return
			}

			var (
				id         = r.Random(p.Rand)
				channelID  = fmt.Sprintf("ch%06d", id)
				messageID  = fmt.Sprintf("%07d", counter)
				messageLog = log.With(zap.String("channel_id", channelID), zap.String("message_id", messageID))
				message    = fmt.Sprintf("Message %s from %s", messageID, r)
				result     = p.Redis.Publish(ctx, channelID, message)
			)
			if err := result.Err(); err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}

				messageLog.Error("Failed to publish message", zap.Error(result.Err()))
				continue
			}

			subscribersCount, err := result.Result()
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				messageLog.Error("Failed to get the amount of subscribers from publish result", zap.Error(err))
				continue
			}

			messageLog.Info(fmt.Sprintf("Message is published to %d clients", subscribersCount))
		}
	}
}
