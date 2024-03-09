package logging

import "github.com/google/uuid"

type Level string

const (
	Info  = "info"
	Debug = "debug"
	Warn  = "warn"
	Error = "error"
)

type Log struct {
	ID    uuid.UUID `json:"id"`
	Level Level     `json:"level"`
	Text  string    `json:"text"`
}
