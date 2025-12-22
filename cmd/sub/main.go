package main

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"

	"github.com/goldobin/microstreams/internal/subsvc"
)

func main() {
	var (
		logger        = zerolog.New(os.Stderr).With().Timestamp().Logger()
		rdb           = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
		svcLogger     = logger.With().Str("component", "subscription-service").Logger()
		svc           = subsvc.Svc{Logger: svcLogger, Redis: rdb}
		handlerLogger = logger.With().Str("component", "http-handler").Logger()
		h             = handler{Logger: handlerLogger, Svc: svc}
		r             = mux.NewRouter()
	)

	r.Handle("/subscription/{channel}", &h)
	if err := http.ListenAndServe(":8080", r); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error().Err(err).Msg("server stopped unexpectedly")
	}
}

type handler struct {
	Logger zerolog.Logger
	Svc    subsvc.Svc
	U      websocket.Upgrader
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name, ok := mux.Vars(r)["channel"]
	if !ok {
		h.writeText(w, http.StatusBadRequest, "Channel is not specified")
		return
	}

	logger := h.Logger.With().Str("channel", name).Logger()

	wsConn, err := h.U.Upgrade(w, r, nil)
	if err != nil {
		h.writeText(w, http.StatusInternalServerError, "WebSocket upgrade failed")
		logger.Error().Err(err).Msg("Failed to upgrade websocket")
		return
	}

	sub := h.Svc.Subscribe(r.Context(), name)
	wsConn.SetCloseHandler(func(code int, text string) error {
		if err := sub.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close subscription")
		}
		return nil
	})

	forwarderLogger := logger.With().Str("component", "redis-consumer").Logger()
	forward(forwarderLogger, wsConn, sub.Messages(context.Background())) // TODO: Use some kind of cancelable context
}

func (h *handler) writeText(w http.ResponseWriter, status int, text string) {
	log := h.Logger

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	written, err := w.Write([]byte(text))
	if err != nil {
		log.Error().Err(err).Msg("Failed to write to websocket")
		return
	}

	if written != len(text) {
		log.Error().
			Int("written_bytes", written).
			Int("total_bytes", len(text)).
			Msg("Failed to write response")
	}
}

func forward(logger zerolog.Logger, dst *websocket.Conn, src <-chan subsvc.Message) {
	for m := range src {
		logger.Print("Message received from redis, sending message to websocket")
		if err := dst.WriteJSON(m); err != nil {
			logger.Error().Err(err).Msg("Failed to write to WebSocket")
			return
		}

		logger.Print("Message is sent to WebSocket")
	}
}
