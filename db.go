package db

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	DEFAULT_MAX_IDLE_CONNS = 10
	DEFAULT_MAX_OPEN_CONNS = 20
	DEFAULT_TIMEOUT        = "90s"
)

type (
	Options struct {
		DriverName     string
		DataSourceName string
		DbName         string
		Timeout        string
		MaxIdleConns   int
		MaxOpenConns   int
	}

	Db struct {
		Options Options
		conn    *sql.DB
	}
)

var connection Db

func CreateConnection(options Options) (*sql.DB, error) {
	connection = Db{Options: options}

	var timeout string
	if options.Timeout != "" {
		timeout = options.Timeout
	} else {
		timeout = DEFAULT_TIMEOUT
	}

	conn, err := sql.Open(options.DriverName, fmt.Sprintf("%s/%s?timeout=%s", options.DataSourceName, options.DbName, timeout))
	if err != nil {
		return nil, errors.New("Error on initializing database connection")
	}

	var (
		maxIdleConns int = DEFAULT_MAX_IDLE_CONNS
		maxOpenConns int = DEFAULT_MAX_OPEN_CONNS
	)

	if options.MaxIdleConns != 0 {
		maxIdleConns = options.MaxIdleConns
	}
	if options.MaxOpenConns != 0 {
		maxOpenConns = options.MaxOpenConns
	}

	conn.SetMaxIdleConns(maxIdleConns)
	conn.SetMaxOpenConns(maxOpenConns)

	connection.conn = conn

	return connection.conn, nil
}

func reconn() (*sql.DB, error) {
	conn, err := sql.Open(connection.Options.DriverName, fmt.Sprintf("%s/%s?timeout=%s", connection.Options.DataSourceName, connection.Options.DbName, connection.Options.Timeout))
	if err != nil {
		return nil, errors.New("Error on initializing database connection")
	}

	conn.SetMaxIdleConns(connection.Options.MaxIdleConns)
	conn.SetMaxOpenConns(connection.Options.MaxOpenConns)

	connection.conn = conn

	return connection.conn, nil
}

func GetConnection() (*sql.DB, error) {
	if err := connection.conn.Ping(); err != nil {
		return reconn()
	}

	return connection.conn, nil
}
