package nats

import (
	"fmt"

	"nats+clickhouse/logging"

	"github.com/nats-io/nats.go"
)

const Subject = "log_subj"

func New() *nats.EncodedConn {
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)

	c.Subscribe(Subject, func(subj, reply string, l *logging.Log) {
		logging.Logger.Info(fmt.Sprintf("Received a logging msg on subject %s! %+v", subj, l))
	})

	return c
}
