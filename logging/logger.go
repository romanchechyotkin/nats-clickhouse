package logging

import (
	"io"
	"log"
	"log/slog"
	"os"
)

var Logger *slog.Logger

func New() {
	file, err := os.OpenFile("logfile.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	multiWriter := io.MultiWriter(file, os.Stdout)

	Logger = slog.New(slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}
