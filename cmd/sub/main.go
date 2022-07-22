package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
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

	ws := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		name, ok := vars["channel"]

		if !ok {
			w.WriteHeader(400)
			_, _ = w.Write([]byte("channel is not specified"))
			return
		}

		wsLog := log.With(
			zap.Namespace("subscription"),
			zap.String("channel", name),
		)

		u := websocket.Upgrader{}
		c, err := u.Upgrade(w, r, nil)

		if err != nil {
			w.WriteHeader(500)
			_, _ = w.Write([]byte("WebSocket upgrade failed"))
			wsLog.Error(
				"failed to upgrade to WebSocket",
				zap.Error(err),
			)
			return
		}

		ctx := r.Context()

		ps := rdb.Subscribe(ctx, name)
		defer func() {
			wsLog.Debug("unsubscribing from redis channel")
			err := ps.Close()
			if err != nil {
				wsLog.Error(
					"failed to unsubscribe from redis channel",
					zap.Error(err),
				)
			}
			log.Debug("successfully unsubscribe from redis channel")
		}()

		wsLog.Debug("subscribed to redis channel")

		go func() {
			consumerLog := wsLog.With(zap.Namespace("redis-consumer"))
			ch := ps.Channel()

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

		consumerLog := wsLog.With(zap.Namespace("ws-consumer"))
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

	router := mux.NewRouter()
	router.HandleFunc("/subscription/{channel}", ws)

	panic(http.ListenAndServe("0.0.0.0:8080", router))
}
