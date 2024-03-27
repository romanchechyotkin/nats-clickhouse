package logging

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

// Buffer is an implementation of a ring buffer
type Buffer struct {
	data          []*Log
	size          int
	head          int
	tail          int
	exportTimeout time.Duration
	nats          *nats.EncodedConn
}

func NewRingBuffer(size int, conn *nats.EncodedConn) *Buffer {
	buf := &Buffer{
		data:          make([]*Log, size),
		size:          size,
		head:          0,
		tail:          -1,
		exportTimeout: 10 * time.Second,
		nats:          conn,
	}

	go buf.exportLogs()

	return buf
}

func (buf *Buffer) Insert(level, text string) {
	data := Log{
		ID:    uuid.New(),
		Level: "",
		Text:  "",
	}

	buf.tail = (buf.tail + 1) % buf.size
	buf.data[buf.tail] = &data

	if buf.head == buf.tail {
		buf.head = (buf.head + 1) % buf.size
	}
}

func (buf *Buffer) Emit() []*Log {
	out := []*Log{}

	for {
		if buf.data[buf.head] != nil {
			out = append(out, buf.data[buf.head])
			buf.data[buf.head] = nil
		}
		if buf.head == buf.tail || buf.tail == 1 {
			break
		}
		buf.head = (buf.head + 1) % buf.size
	}

	return out
}

func (buf *Buffer) exportLogs() {
	ticker := time.NewTicker(buf.exportTimeout)

	for {
		select {
		case <-ticker.C:
			file, err := os.Open("logfile.log")
			if err != nil {
				Logger.Error(err.Error())
				continue
			}

			scanner := bufio.NewReader(file)

			var logMap map[string]any

			for i := 0; i < buf.size; i++ {
				line, b, err := scanner.ReadLine()
				log.Println("from file", string(line), b, err)

				err = json.Unmarshal(line, &logMap)
				if err != nil {
					Logger.Error(err.Error())
				}

				buf.Insert(logMap["level"].(string), logMap["msg"].(string))
			}

			//logs := buf.Emit()
			//Logger.Info("logs", slog.Int("length", len(logs)))
		}
	}
}
