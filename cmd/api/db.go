package main

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	Open(config) (*gorm.DB, error)
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

func (m MySQLDB) Open(cfg config) (*gorm.DB, error) {
	c := MySQLDB{
		database: os.Getenv("RDB_DBNAME"),
		username: os.Getenv("RDB_USERNAME"),
		pwd:      os.Getenv("RDB_PASSWORD"),
		port:     os.Getenv("RDB_PORT"),
		host:     os.Getenv("RDB_HOST"),
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.username, c.pwd, c.host, c.port, c.database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (m PSQLDB) Open(cfg config) (*gorm.DB, error) {
	c := PSQLDB{
		database: os.Getenv("DB_DBNAME"),
		username: os.Getenv("DB_USERNAME"),
		pwd:      os.Getenv("DB_PASSWORD"),
		port:     os.Getenv("DB_PORT"),
		host:     os.Getenv("DB_HOST"),
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", c.host, c.username, c.pwd, c.database, c.port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
