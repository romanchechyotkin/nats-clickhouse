package main

import (
	"log/slog"
	"net/http"

	"nats+clickhouse/clickhouse"
	"nats+clickhouse/logging"
	"nats+clickhouse/nats"
)

func main() {
	logging.New()
	clickhouse.New()
	natsConn := nats.New()

	logging.Logger.Info("running server", slog.Int("port", 8080))

	logging.NewRingBuffer(10, natsConn)

	http.HandleFunc("/click", func(w http.ResponseWriter, req *http.Request) {
		logging.Logger.Info("got request", slog.String("uri", req.RequestURI), slog.String("method", req.Method))

		_, _ = w.Write([]byte("clickhouse"))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		logging.Logger.Error(err.Error())
	}
}
