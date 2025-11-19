package main

import (
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func main() {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	defer func(log *zap.Logger) {
		if err := log.Sync(); err != nil {
			panic(err)
		}
	}(log)

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	channelHandler := newChannelHandler(log.Named("channel-handler"), rdb)

	r := mux.NewRouter()
	r.HandleFunc("/subscription/{channel}", channelHandler)

	if err := http.ListenAndServe(":8080", r); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("server stopped unexpectedly", zap.Error(err))
	}
}

func newChannelHandler(log *zap.Logger, rdb *redis.Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name, ok := vars["channel"]

		if !ok {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("channel is not specified"))
			return
		}

		log = log.With(zap.String("channel", name))

		u := websocket.Upgrader{}
		c, err := u.Upgrade(w, r, nil)
		if err != nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("WebSocket upgrade failed"))
			log.Error(
				"failed to upgrade to WebSocket",
				zap.Error(err),
			)
			return
		}

		s := rdb.Subscribe(r.Context(), name)
		defer func() {
			log.Debug("Unsubscribing from redis channel")
			if err := s.Close(); err != nil {
				log.Error("Failed to unsubscribe from redis channel", zap.Error(err))
			}
			log.Debug("successfully unsubscribe from redis channel")
		}()

		log.Debug("subscribed to redis channel")

		go func() {
			consumerLog := log.With(zap.Namespace("redis-consumer"))
			ch := s.Channel()

			for {
				msg, ok := <-ch

				if !ok {
					consumerLog.Debug("redis channel is closed")
					return
				}

				consumerLog.Debug("message received from redis, sending message to websocket")

				if err = c.WriteJSON(msg); err != nil {
					consumerLog.Error("failed to write to websocket", zap.Error(err))
					return
				}

				consumerLog.Debug("message is sent to WS")
			}
		}()

		consumerLog := log.With(zap.Namespace("ws-consumer"))
		for {
			_, _, err := c.NextReader()
			if err != nil {
				consumerLog.Error(
					"failed to read from websocket",
					zap.Error(err),
				)
				break
			}

			consumerLog.Debug("message is received from websocket, ignoring")
		}
	}
}
