package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func main() {
	var (
		log, logSync = newLogger()
		rdb          = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
		h            = handler{Log: log.Named("handler"), Redis: rdb}
		r            = mux.NewRouter()
	)
	defer logSync()

	r.Handle("/subscription/{channel}", &h)
	if err := http.ListenAndServe(":8080", r); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("server stopped unexpectedly", zap.Error(err))
	}
}

type message struct {
	Text string `json:"text"`
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

type handler struct {
	Log   *zap.Logger
	Redis *redis.Client
	U     websocket.Upgrader
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name, ok := mux.Vars(r)["channel"]
	if !ok {
		h.writeText(w, http.StatusBadRequest, "Channel is not specified")
		return
	}

	log := h.Log.With(zap.String("channel", name))

	wsConn, err := h.U.Upgrade(w, r, nil)
	if err != nil {
		h.writeText(w, http.StatusInternalServerError, "WebSocket upgrade failed")
		log.Error(
			"failed to upgrade to WebSocket",
			zap.Error(err),
		)
		return
	}

	s := h.Redis.Subscribe(r.Context(), name)
	ms := convert(toMessage, s.Channel())
	forward(log.With(zap.Namespace("redis-consumer")), wsConn, ms)
}

func (h *handler) writeText(w http.ResponseWriter, status int, text string) {
	log := h.Log

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	written, err := w.Write([]byte(text))
	if err != nil {
		log.Error("Failed to write to websocket", zap.Error(err))
		return
	}

	if written != len(text) {
		log.Error(
			"Failed to write response",
			zap.Int("actually_written", written),
			zap.Int("should_be_written", len(text)),
		)
	}
}

func convert(conv func(*message, *redis.Message), src <-chan *redis.Message) <-chan message {
	out := make(chan message)
	go func() {
		defer close(out)
		for rm := range src {
			var m message
			conv(&m, rm)
			out <- m
		}
	}()

	return out
}

func forward(log *zap.Logger, dst *websocket.Conn, src <-chan message) {
	for m := range src {
		log.Debug("message received from redis, sending message to websocket")
		if err := dst.WriteJSON(m); err != nil {
			log.Error("Failed to write to WebSocket", zap.Error(err))
			return
		}

		log.Debug("message is sent to WebSocket")
	}
}

func toMessage(dst *message, src *redis.Message) {
	dst.Text = src.Payload
}
