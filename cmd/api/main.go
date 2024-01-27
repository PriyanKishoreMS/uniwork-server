package main

import (
	"flag"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type config struct {
	port int
	env  string
}

type application struct {
	config config
}

func main() {
	cfg := config{}

	flag.IntVar(&cfg.port, "port", 3000, "Port number")
	flag.StringVar(&cfg.env, "env", "development", "Environment")
	flag.Parse()
	log.SetHeader("${time_rfc3339} ${level}")

	dbType := PSQLDB{}
	db, err := dbType.Open(cfg)
	if err != nil {
		log.Fatal("error in opening db", err)
	}

	log.Info("Database connection established")
	dbname := db.Migrator().CurrentDatabase()
	log.Info("dbname: ", dbname)

	app := &application{
		config: cfg,
	}

	e := echo.New()
	e.HideBanner = true
	app.routes(e)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", app.config.port)))
}
