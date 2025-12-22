package subsvc

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

type Svc struct {
	Logger zerolog.Logger
	Redis  *redis.Client
}

type Message struct {
	Text string `json:"text"`
}

type ConvFunc func(*Message, *redis.Message)

type Subscription struct {
	Logger zerolog.Logger
	PubSub *redis.PubSub
}

func (s *Subscription) Messages(ctx context.Context) <-chan Message {
	return s.ConvertChan(ctx, s.PubSub.Channel())
}

func (s *Subscription) Close() error {
	return s.PubSub.Close()
}

func (s *Svc) Subscribe(ctx context.Context, channel string) Subscription {
	return Subscription{
		Logger: s.Logger.With().Str("channel", channel).Logger(),
		PubSub: s.Redis.Subscribe(ctx, channel),
	}
}

func (s *Subscription) ConvertChan(ctx context.Context, src <-chan *redis.Message) <-chan Message {
	out := make(chan Message)
	go func() {
		defer func() {
			close(out)
			s.Logger.Print("Out channel closed")
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case redisMessage, ok := <-src:
				if !ok {
					return
				}

				s.Logger.Println("Message received from Redis")
				var m Message
				redisMessageToMessage(&m, redisMessage)

				select {
				case <-ctx.Done():
					return
				case out <- m:
					s.Logger.Println("Message converted and sent to out channel")
				}
			}
		}
	}()

	return out
}

func redisMessageToMessage(dst *Message, src *redis.Message) {
	dst.Text = src.Payload
}
