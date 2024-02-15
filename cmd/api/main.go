package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"firebase.google.com/go/messaging"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	// fb       data.FirebaseUtils
	fcmClient *messaging.Client
	awsS3     *data.S3
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

	dbType := PSQLDB{}
	db, err := dbType.Open(cfg)
	if err != nil {
		log.Fatalf("error in opening db; %v", err)
	}
	defer db.Close()

	validate = *validator.New()

	firebase, err := data.NewFirebaseIntegration()
	if err != nil {
		log.Fatalf("error in firebase util: %v", err)
	}

	fcmClient, err := firebase.App.Messaging(context.Background())
	if err != nil {
		log.Fatalf("error in initializing fcmclient: %v", err)
	}

	s3cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("error in loading aws config: %v", err)
	}
	client := s3.NewFromConfig(s3cfg)

	app := &application{
		config:    cfg,
		models:    data.NewModel(db),
		validate:  validate,
		fcmClient: fcmClient,
		awsS3:     data.NewS3(client),
	}
	e := app.routes()
	e.Server.ReadTimeout = time.Second * 10
	e.Server.WriteTimeout = time.Second * 20
	e.Server.IdleTimeout = time.Minute
	e.HideBanner = true

	log.Info("Server starting on port: ", cfg.port)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.port)))
}
