package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Database interface {
	Open(config) (*sql.DB, error)
}

type MySQLDB struct {
	database string
	username string
	pwd      string
	port     string
	host     string
}

type PSQLDB struct {
	database string
	username string
	pwd      string
	port     string
	host     string
}

func (m MySQLDB) Open(cfg config) (*sql.DB, error) {
	c := MySQLDB{
		database: os.Getenv("RDB_DBNAME"),
		username: os.Getenv("RDB_USERNAME"),
		pwd:      os.Getenv("RDB_PASSWORD"),
		port:     os.Getenv("RDB_PORT"),
		host:     os.Getenv("RDB_HOST"),
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.username, c.pwd, c.host, c.port, c.database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (m PSQLDB) Open(cfg config) (*sql.DB, error) {
	c := PSQLDB{
		database: os.Getenv("DB_DATABASE"),
		username: os.Getenv("DB_USERNAME"),
		pwd:      os.Getenv("DB_PASSWORD"),
		port:     os.Getenv("DB_PORT"),
		host:     os.Getenv("DB_HOST"),
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", c.host, c.username, c.pwd, c.database, c.port)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
