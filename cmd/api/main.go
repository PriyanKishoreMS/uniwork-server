package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
	"github.com/priyankishorems/uniwork-server/internal/data"
)

type Config struct {
	port int
	env  string
	jwt  struct {
		secret string
		issuer string
	}
	limiter struct {
		rps     int
		burst   int
		enabled bool
	}
}

type application struct {
	config   Config
	models   data.Models
	validate validator.Validate
	// wg       sync.WaitGroup
}

var validate validator.Validate

func main() {
	cfg := Config{}

	flag.IntVar(&cfg.port, "port", 3000, "Port number")
	flag.StringVar(&cfg.env, "env", "development", "Environment")

	flag.StringVar(&cfg.jwt.secret, "jwt-secret", os.Getenv("JWT_SECRET"), "JWT secret")
	flag.StringVar(&cfg.jwt.issuer, "jwt-issuer", os.Getenv("JWT_ISSUER"), "JWT issuer")

	flag.IntVar(&cfg.limiter.rps, "limiter-rps", 10, "Rate limiter max requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 10, "Rate limiter max burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Rate limiter enabled")

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
	e.Server.ReadTimeout = time.Second * 10
	e.Server.WriteTimeout = time.Second * 20
	e.Server.IdleTimeout = time.Minute
	e.HideBanner = true

	log.Info("Server starting on port: ", cfg.port)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.port)))
}
