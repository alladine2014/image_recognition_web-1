package storage

import (
	"database/sql"
	"fmt"
	"github.com/cgCodeLife/image_recognition_web/config"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

type Config struct {
	DSN         string
	Active      int
	Idle        int
	IdleTimeout time.Duration
}

var (
	storage *Storage
)

type Storage struct {
	c  *Config
	db *sql.DB
}

func new(c *Config) *Storage {
	return &Storage{
		c:  c,
		db: newMySql(c),
	}
}

func newMySql(c *Config) *sql.DB {
	db, err := open(c)
	if err != nil {
		panic(err)
	}
	return db
}

func open(c *Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", c.DSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("sql.Open() error(%v)", err))
		return nil, err
	}
	db.SetMaxOpenConns(c.Active)
	db.SetMaxIdleConns(c.Idle)
	db.SetConnMaxLifetime(c.IdleTimeout)
	return db, nil
}

func Init() {
	conf := &Config{
		DSN:         config.GetMysql().GetDSN(),
		Active:      config.GetMysql().GetActive(),
		Idle:        config.GetMysql().GetIdle(),
		IdleTimeout: config.GetMysql().GetIdleTimeout(),
	}
	storage = new(conf)
}

func (s *Storage) Close() {
	s.db.Close()
}
