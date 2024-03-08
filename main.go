package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/google/uuid"
	_ "github.com/mailru/go-clickhouse"
	"github.com/nats-io/nats.go"
)

const subject = "log_subj"

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

func main() {
	clickhouseCfg := ClickhouseConfig{
		Host:     "localhost",
		Port:     "8123",
		User:     "clickhouse",
		Password: "8123",
	}

	clickhouseClient := newClickhouseClient(&clickhouseCfg)
	_ = clickhouseClient
	runMigration(clickhouseClient)

	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer c.Close()

	ch := make(chan struct{})

	c.Subscribe(subject, func(subj, reply string, l *Log) {
		fmt.Printf("Received a log msg on subject %s! %+v\n", subj, l)
		ch <- struct{}{}
	})

	me := &Log{ID: uuid.New(), Level: Info, Text: "robert lox"}
	c.Publish(subject, me)

	<-ch
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

//go:embed migrations/*.sql
var migrations embed.FS

func runMigration(clickhouseClient *sql.DB) {
	dir, err := migrations.ReadDir("migrations")
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range dir {
		content, err := migrations.ReadFile("migrations/" + entry.Name())
		if err != nil {
			log.Fatal(err)
		}

		log.Println(string(content))

		exec, err := clickhouseClient.Exec(string(content))
		if err != nil {
			log.Fatal(err)
		}

		log.Println(exec.RowsAffected())
	}

}
