package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
	"github.com/priyankishorems/uniwork-server/internal/data"
)

type Config struct {
	port int
	env  string
}

type application struct {
	config   Config
	models   data.Models
	validate validator.Validate
}

var validate validator.Validate

func main() {
	cfg := Config{}

	flag.IntVar(&cfg.port, "port", 3000, "Port number")
	flag.StringVar(&cfg.env, "env", "development", "Environment")
	flag.Parse()
	log.SetHeader("${time_rfc3339} ${level}")

	dbType := MySQLDB{}
	db, err := dbType.Open(cfg)
	if err != nil {
		log.Fatal("error in opening db", err)
	}
	defer db.Close()

	validate = *validator.New()

	app := &application{
		config:   cfg,
		models:   data.NewModel(db),
		validate: validate,
	}
	e := app.routes()
	e.Server.ReadHeaderTimeout = time.Second * 10
	e.Server.WriteTimeout = time.Second * 20
	e.Server.IdleTimeout = time.Minute
	e.HideBanner = true

	log.Info("Server starting on port: ", cfg.port)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.port)))
}
