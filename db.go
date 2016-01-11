package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	DEFAULT_MAX_IDLE_CONNS = 100
	DEFAULT_MAX_OPEN_CONNS = 100
	DEFAULT_TIMEOUT        = "90s"
	TIMEOUT_CONNECTION     = time.Second * 20
)

type Options struct {
	DriverName     string
	DataSourceName string
	DbName         string
	Timeout        string
	MaxIdleConns   int
	MaxOpenConns   int
}

type Db struct {
	Options   Options
	conn      *sql.DB
	last_time time.Time
	Timer     *time.Timer
}

var connection Db

func CreateConnection(options Options) (*sql.DB, error) {
	connection = Db{Options: options}

	var timeout string = DEFAULT_TIMEOUT
	if options.Timeout != "" {
		timeout = options.Timeout
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

/**
	Соединение с базой
**/
/*
func (o Options) CreateConnection_() (*sql.DB, error) {
	connection.Options = o

	var timeout string = DEFAULT_TIMEOUT
	if o.Timeout != "" {
		timeout = o.Timeout
	}

	conn, err := sql.Open(o.DriverName, o.DataSourceName+"/"+o.DbName+"?timeout="+timeout)
	if err != nil {
		log.Fatalf("Error on initializing database connection: %s", err.Error())
		return nil, errors.New("Error on initializing database connection")
	}

	var maxIdleConns int = DEFAULT_MAX_IDLE_CONNS
	var maxOpenConns int = DEFAULT_MAX_OPEN_CONNS

	if o.MaxIdleConns != 0 {
		maxIdleConns = o.MaxIdleConns
	}
	if o.MaxOpenConns != 0 {
		maxOpenConns = o.MaxOpenConns
	}

	conn.SetMaxIdleConns(maxIdleConns)
	conn.SetMaxOpenConns(maxOpenConns)

	connection.Connection = conn

	if connection.Timer != nil {
		connection.Timer.Stop()
	}
	connection.Timer = time.NewTimer(TIMEOUT_CONNECTION)

	go func() {
		<-connection.Timer.C
		connection.Connection.Close()
		connection.Connection = nil
	}()

	return conn, nil
}
*/

/*
func Connection() (*sql.DB, error) {
	var err error

	if connection.Connection == nil {
		connection.Connection, err = connection.Options.CreateConnection()
		if err != nil {
			connection.Timer.Stop()
			connection.Connection = nil

			return nil, err
		}
	} else {
		if err = connection.Connection.Ping(); err != nil {
			connection.Timer.Stop()
			connection.Connection = nil

			return nil, err
		}
	}

	connection.Timer.Reset(TIMEOUT_CONNECTION)

	return connection.Connection, nil
}
*/
