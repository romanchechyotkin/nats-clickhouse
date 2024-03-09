package clickhouse

import (
	"database/sql"
	"embed"
	"fmt"

	"nats+clickhouse/logging"
)

func New() {
	cfg := Config{
		Host:     "localhost",
		Port:     "8123",
		User:     "clickhouse",
		Password: "8123",
	}

	connStr := fmt.Sprintf("http://%s:%s/?user=%s&password=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password)
	driver := "clickhouse"
	connect, err := sql.Open(driver, connStr)
	if err != nil {
		logging.Logger.Error(err.Error())
	}

	if err := connect.Ping(); err != nil {
		logging.Logger.Error(err.Error())
	}

	var dbName string
	err = connect.QueryRow("SELECT currentDatabase()").Scan(&dbName)
	if err != nil {
		logging.Logger.Error(err.Error())
	}

	logging.Logger.Info("clickhouse current database:", dbName)

	RunMigration(connect)
}

//go:embed migrations/*.sql
var migrations embed.FS

func RunMigration(clickhouseClient *sql.DB) {
	dir, err := migrations.ReadDir("migrations")
	if err != nil {
		logging.Logger.Error(err.Error())
	}

	for _, entry := range dir {
		content, err := migrations.ReadFile("migrations/" + entry.Name())
		if err != nil {
			logging.Logger.Error(err.Error())
		}

		_, err = clickhouseClient.Exec(string(content))
		if err != nil {
			logging.Logger.Error(err.Error())
		}
	}

}
