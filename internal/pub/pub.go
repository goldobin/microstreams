package pub

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"

	"github.com/goldobin/microstreams/internal/ranges"
	"github.com/goldobin/microstreams/internal/rate"
)

type Pub struct {
	Logger zerolog.Logger
	Redis  *redis.Client
	Rand   *rand.Rand
	Rate   rate.Rate
}

type contextKey struct{}

var shardContextKey contextKey

func (p Pub) Publish(ctx context.Context, shard ranges.Int) {
	ctx = context.WithValue(ctx, shardContextKey, shard)
	time.Sleep(p.Rate.Random(p.Rand))
	var (
		counter    = 0
		logger     = p.Logger.With().Ctx(ctx).Logger()
		ticker     = time.NewTicker(p.Rate.Interval())
		stopReason = "undefined"
	)
	defer ticker.Stop()

	logger.Info().Msg("Start publishing...")
	defer func() {
		logger.Info().Str("stop_reason", stopReason).Msg("Publishing stopped")
	}()

	for {
		counter++
		select {
		case <-ctx.Done():
			stopReason = ctx.Err().Error()
			return
		case _, ok := <-ticker.C:
			if !ok {
				stopReason = "ticker stopped"
				return
			}

			var (
				id        = shard.Random(p.Rand)
				channelID = fmt.Sprintf("ch%06d", id)
				messageID = fmt.Sprintf("%07d", counter)
				msgLogger = logger.
						With().
						Str("channel_id", channelID).
						Str("message_id", messageID).
						Logger()
				msg    = fmt.Sprintf("Message %s from %s", messageID, shard)
				result = p.Redis.Publish(ctx, channelID, msg)
			)
			if err := result.Err(); err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				msgLogger.
					Error().
					Err(result.Err()).
					Msg("Failed to publish message")
				continue
			}

			subscribersCount, err := result.Result()
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				msgLogger.
					Error().
					Err(err).
					Msg("Failed to get the amount of subscribers from publish result")
				continue
			}

			msgLogger.Info().Int64("subscriber_count", subscribersCount).Msg("Message published")
		}
	}
}
