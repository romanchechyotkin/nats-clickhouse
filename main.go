package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mailru/go-clickhouse"
	"github.com/nats-io/nats.go"
)

func main() {
	clickhouseCfg := ClickhouseConfig{
		Host:     "localhost",
		Port:     "8123",
		User:     "clickhouse",
		Password: "8123",
	}

	clickhouseClient := newClickhouseClient(&clickhouseCfg)
	_ = clickhouseClient

	nc, _ := nats.Connect(nats.DefaultURL)

	// Channel Subscriber
	ch := make(chan *nats.Msg, 64)
	_, err := nc.ChanSubscribe("foo", ch)
	if err != nil {
		log.Println(err)
	}

	// Simple Publisher
	err = nc.Publish("foo", []byte("Hello World"))
	if err != nil {
		log.Println(err)
	}

	msg := <-ch
	log.Println(string(msg.Data))
}

type ClickhouseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

func newClickhouseClient(cfg *ClickhouseConfig) *sql.DB {
	connStr := fmt.Sprintf("http://%s:%s/?user=%s&password=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password)
	driver := "clickhouse"
	connect, err := sql.Open(driver, connStr)
	if err != nil {
		log.Fatalf("Open >> %v", err)
	}

	if err := connect.Ping(); err != nil {
		log.Fatalf("Ping >> %v", err)
	}

	var dbName string
	err = connect.QueryRow("SELECT currentDatabase()").Scan(&dbName)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("clickhouse current database:", dbName)

	return connect
}
