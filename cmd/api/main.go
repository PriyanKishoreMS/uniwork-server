package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/priyankishorems/uniwork-server/internal/data"
	"github.com/priyankishorems/uniwork-server/internal/database"
)

type Config struct {
	port int
	env  string
}

type application struct {
	config Config
	models data.Models
}

func main() {
	cfg := Config{}

	flag.IntVar(&cfg.port, "port", 3000, "Port number")
	flag.StringVar(&cfg.env, "env", "development", "Environment")
	flag.Parse()
	log.SetHeader("${time_rfc3339} ${level}")

	dbType := database.MySQLDB{}
	db, err := dbType.Open()
	if err != nil {
		log.Fatal("error in opening db", err)
	}
	defer db.Close()

	log.Info("Database connection established")

	app := &application{
		config: cfg,
		models: data.NewModel(db),
	}
	e := app.routes()
	e.Server.ReadHeaderTimeout = time.Second * 10
	e.Server.WriteTimeout = time.Second * 20
	e.Server.IdleTimeout = time.Minute
	e.HideBanner = true

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.port)))
}
