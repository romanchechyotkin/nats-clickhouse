package logging

import (
	"github.com/google/uuid"
)

// Buffer is an implementation of a ring buffer
type Buffer struct {
	data []*Log
	size int
	head int
	tail int
}

func NewRingBuffer(size int) *Buffer {
	return &Buffer{
		data: make([]*Log, size),
		size: size,
		head: 0,
		tail: -1,
	}
}

func (buf *Buffer) Insert() {
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
