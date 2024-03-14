package main

import (
	"context"
	"time"

	_ "go.uber.org/automaxprocs"

	"booking-schedule/internal/pkg/auth"
	"flag"
	"log"
)

var pathConfig, pathCert, pathKey string

func init() {
	flag.StringVar(&pathConfig, "config", ".configs/auth_config.yml", "path to config file")
	flag.StringVar(&pathCert, "certfile", "cert.pem", "certificate PEM file")
	flag.StringVar(&pathKey, "keyfile", "key.pem", "key PEM file")
	time.Local = time.UTC
}

func main() {
	flag.Parse()

	ctx := context.Background()

	app, err := auth.NewApp(ctx, pathConfig, pathCert, pathKey)
	if err != nil {
		log.Fatalf("failed to create auth-api app object:%s\n", err.Error())
	}

	err = app.Run()
	if err != nil {
		log.Fatalf("failed to run auth-api app: %s", err.Error())
	}
}
