package main

import (
	"log/slog"
	"net/http"

	"nats+clickhouse/clickhouse"
	"nats+clickhouse/logging"
	"nats+clickhouse/nats"

	"github.com/google/uuid"
	_ "github.com/mailru/go-clickhouse"
)

func main() {
	logging.New()
	clickhouse.New()
	natsConn := nats.New()

	logging.Logger.Info("running server", slog.Int("port", 8080))

	http.HandleFunc("/click", func(w http.ResponseWriter, req *http.Request) {
		logging.Logger.Info("got request", slog.String("uri", req.RequestURI), slog.String("method", req.Method))

		go func() {
			me := &logging.Log{ID: uuid.New(), Level: logging.Info, Text: "robert lox"}
			err := natsConn.Publish(nats.Subject, me)
			if err != nil {
				logging.Logger.Error("failed to publish message", err)
			}
		}()

		_, _ = w.Write([]byte("clickhouse"))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		logging.Logger.Error(err.Error())
	}
}
